package honey

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"honey_back/internal/models"
	"honey_back/internal/utils/res"
)

type Service struct {
	db  *gorm.DB
	ops ContainerOps
}

func NewService(db *gorm.DB, ops ContainerOps) *Service {
	return &Service{db: db, ops: ops}
}

// ---------- 镜像模板 ----------

func (s *Service) ListImages(ctx context.Context) ([]models.HoneypotImage, error) {
	var list []models.HoneypotImage
	err := s.db.WithContext(ctx).Order("id ASC").Find(&list).Error
	return list, err
}

// ---------- 部署实例 ----------

type CreateDeploymentReq struct {
	Name    string `json:"name" binding:"required,max=100"`
	ImageID int64  `json:"image_id" binding:"required"`
}

// DeploymentView 供接口返回的视图
type DeploymentView struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	ImageID       int64     `json:"image_id"`
	ImageName     string    `json:"image_name"`
	ImageType     string    `json:"image_type"`
	Status        string    `json:"status"`
	HostPort      int       `json:"host_port"`
	ContainerPort int       `json:"container_port"`
	ContainerID   string    `json:"container_id"`
	CreatedAt     time.Time `json:"created_at"`
}

func (s *Service) loadImages(ctx context.Context, ids []int64) (map[int64]models.HoneypotImage, error) {
	out := make(map[int64]models.HoneypotImage, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	var imgs []models.HoneypotImage
	if err := s.db.WithContext(ctx).Where("id IN ?", ids).Find(&imgs).Error; err != nil {
		return nil, err
	}
	for _, im := range imgs {
		out[im.ID] = im
	}
	return out, nil
}

func (s *Service) toView(d *models.Deployment, imgs map[int64]models.HoneypotImage) DeploymentView {
	v := DeploymentView{
		ID:            d.ID,
		Name:          d.Name,
		ImageID:       d.ImageID,
		Status:        d.Status,
		HostPort:      d.HostPort,
		ContainerPort: d.ContainerPort,
		ContainerID:   d.ContainerID,
		CreatedAt:     d.CreatedAt,
	}
	if im, ok := imgs[d.ImageID]; ok {
		v.ImageName = im.Name
		v.ImageType = im.Type
	}
	return v
}

func (s *Service) getOwned(ctx context.Context, userID, id int64) (*models.Deployment, error) {
	var d models.Deployment
	err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&d).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, res.E(res.NotFound, "部署不存在")
	}
	return &d, err
}

// List GET 当前用户的部署
func (s *Service) List(ctx context.Context, userID int64) ([]DeploymentView, error) {
	var ds []models.Deployment
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("id DESC").Find(&ds).Error; err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(ds))
	for i := range ds {
		ids = append(ids, ds[i].ImageID)
	}
	imgs, err := s.loadImages(ctx, ids)
	if err != nil {
		return nil, err
	}
	views := make([]DeploymentView, 0, len(ds))
	for i := range ds {
		views = append(views, s.toView(&ds[i], imgs))
	}
	return views, nil
}

// Get 单个部署
func (s *Service) Get(ctx context.Context, userID, id int64) (DeploymentView, error) {
	d, err := s.getOwned(ctx, userID, id)
	if err != nil {
		return DeploymentView{}, err
	}
	imgs, err := s.loadImages(ctx, []int64{d.ImageID})
	if err != nil {
		return DeploymentView{}, err
	}
	return s.toView(d, imgs), nil
}

// Create 创建部署：插行(creating) → 拉镜像 → 建容器 → 启动 → 置 running。
// 任一步失败：尽力删除已建容器，行保留为 error 供排障。
func (s *Service) Create(ctx context.Context, userID int64, req *CreateDeploymentReq) (DeploymentView, error) {
	var img models.HoneypotImage
	if err := s.db.WithContext(ctx).First(&img, req.ImageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return DeploymentView{}, res.E(res.ParamCode, "镜像不存在")
		}
		return DeploymentView{}, err
	}

	hostPort, err := s.allocateHostPort(ctx, img.DefaultPort)
	if err != nil {
		return DeploymentView{}, err
	}

	d := models.Deployment{
		Name:          req.Name,
		UserID:        userID,
		ImageID:       img.ID,
		Status:        models.DeploymentCreating,
		HostPort:      hostPort,
		ContainerPort: img.DefaultPort,
	}
	if err := s.db.WithContext(ctx).Create(&d).Error; err != nil {
		return DeploymentView{}, err
	}

	containerName := fmt.Sprintf("hp-%d", d.ID)
	if err := s.ops.EnsureImage(ctx, img.DockerImage); err != nil {
		s.markError(ctx, d.ID)
		return DeploymentView{}, res.E(res.SystemError, fmt.Sprintf("镜像准备失败：%v", err))
	}
	cid, err := s.ops.CreateContainer(ctx, containerName, img.DockerImage, hostPort, img.DefaultPort)
	if err != nil {
		s.markError(ctx, d.ID)
		return DeploymentView{}, res.E(res.SystemError, fmt.Sprintf("容器创建失败：%v", err))
	}
	if err := s.ops.StartContainer(ctx, cid); err != nil {
		_ = s.ops.RemoveContainer(ctx, cid, true)
		s.markError(ctx, d.ID)
		return DeploymentView{}, res.E(res.SystemError, fmt.Sprintf("容器启动失败：%v", err))
	}

	now := time.Now()
	s.db.Model(&d).Updates(map[string]interface{}{
		"container_id": cid,
		"status":       models.DeploymentRunning,
		"updated_at":   now,
	})
	d.ContainerID = cid
	d.Status = models.DeploymentRunning
	return s.toView(&d, map[int64]models.HoneypotImage{img.ID: img}), nil
}

// allocateHostPort 从容器端口起向上找 DB 中未占用的宿主机端口。
func (s *Service) allocateHostPort(ctx context.Context, containerPort int) (int, error) {
	for p := containerPort; p < containerPort+1000; p++ {
		var cnt int64
		if err := s.db.WithContext(ctx).Model(&models.Deployment{}).
			Where("host_port = ?", p).Count(&cnt).Error; err != nil {
			return 0, err
		}
		if cnt == 0 {
			return p, nil
		}
	}
	return 0, res.E(res.ParamCode, "无可用宿主机端口")
}

func (s *Service) markError(ctx context.Context, id int64) {
	s.db.Model(&models.Deployment{}).Where("id = ?", id).Update("status", models.DeploymentError)
}

// Start 启动已停止的实例
func (s *Service) Start(ctx context.Context, userID, id int64) error {
	d, err := s.getOwned(ctx, userID, id)
	if err != nil {
		return err
	}
	if d.Status == models.DeploymentRunning {
		return res.E(res.ParamCode, "实例已在运行")
	}
	if d.ContainerID == "" {
		return res.E(res.SystemError, "容器未创建")
	}
	if err := s.ops.StartContainer(ctx, d.ContainerID); err != nil {
		s.markError(ctx, d.ID)
		return res.E(res.SystemError, fmt.Sprintf("启动失败：%v", err))
	}
	return s.db.Model(&models.Deployment{}).Where("id = ?", d.ID).Update("status", models.DeploymentRunning).Error
}

// Stop 停止运行中的实例
func (s *Service) Stop(ctx context.Context, userID, id int64) error {
	d, err := s.getOwned(ctx, userID, id)
	if err != nil {
		return err
	}
	if d.Status != models.DeploymentRunning {
		return res.E(res.ParamCode, "实例不在运行状态")
	}
	if err := s.ops.StopContainer(ctx, d.ContainerID); err != nil {
		s.markError(ctx, d.ID)
		return res.E(res.SystemError, fmt.Sprintf("停止失败：%v", err))
	}
	return s.db.Model(&models.Deployment{}).Where("id = ?", d.ID).Update("status", models.DeploymentStopped).Error
}

// Delete 删除部署：尽力删除容器（容器可能已不存在），再删 DB 行（v1 硬删除）。
func (s *Service) Delete(ctx context.Context, userID, id int64) error {
	d, err := s.getOwned(ctx, userID, id)
	if err != nil {
		return err
	}
	if d.ContainerID != "" {
		_ = s.ops.RemoveContainer(ctx, d.ContainerID, true)
	}
	return s.db.WithContext(ctx).Delete(&models.Deployment{}, d.ID).Error
}
