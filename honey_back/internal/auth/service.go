package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"honey_back/internal/config"
	"honey_back/internal/models"
	"honey_back/internal/utils"
	"honey_back/internal/utils/res"
)

// Service 会话与登录业务。业务方法只依赖 context.Context。
type Service struct {
	db      *gorm.DB
	rdb     *redis.Client
	session config.Session
}

func NewService(db *gorm.DB, rdb *redis.Client, session config.Session) *Service {
	return &Service{db: db, rdb: rdb, session: session}
}

func (s *Service) ttl() time.Duration {
	return time.Duration(s.session.ExpireHours) * time.Hour
}

// Login 校验账户密码，返回登录态用户（不含敏感字段）。
func (s *Service) Login(ctx context.Context, account, password string) (*models.Principal, error) {
	var u models.User
	err := s.db.WithContext(ctx).Where("user_account = ?", account).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, res.E(res.BadLogin, "账户或密码错误")
	}
	if err != nil {
		return nil, err
	}
	if !utils.ParsePassword(password, u.UserPassword) {
		return nil, res.E(res.BadLogin, "账户或密码错误")
	}

	now := time.Now()
	_ = s.db.Model(&u).Update("last_login_at", now).Error
	return principalOf(&u), nil
}

// CreateSession 生成随机 session id 并把登录态写入 redis。
func (s *Service) CreateSession(ctx context.Context, p *models.Principal) (string, error) {
	raw, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	sid := utils.GenerateSession()
	if err := s.rdb.Set(ctx, s.session.RedisPrefix+sid, raw, s.ttl()).Err(); err != nil {
		return "", err
	}
	return sid, nil
}

// CurrentSession 按 session id 取登录态。
func (s *Service) CurrentSession(ctx context.Context, sid string) (*models.Principal, error) {
	raw, err := s.rdb.Get(ctx, s.session.RedisPrefix+sid).Result()
	if errors.Is(err, redis.Nil) {
		return nil, res.E(res.UserNotLogin, "登录已过期，请重新登录")
	}
	if err != nil {
		return nil, err
	}
	var p models.Principal
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Service) DestroySession(ctx context.Context, sid string) {
	if sid != "" {
		s.rdb.Del(ctx, s.session.RedisPrefix+sid)
	}
}

// BootstrapAdmin 无任何用户时创建一个默认管理员（dev 便利，密码应尽快修改）。
func (s *Service) BootstrapAdmin(ctx context.Context, cfg config.Admin) error {
	var cnt int64
	if err := s.db.WithContext(ctx).Model(&models.User{}).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		return nil
	}
	hash, err := utils.HashPassword(cfg.Password)
	if err != nil {
		return err
	}
	u := models.User{
		UserAccount:  cfg.Account,
		UserPassword: hash,
		UserName:     cfg.UserName,
		UserRole:     models.RoleAdmin,
	}
	if err := s.db.WithContext(ctx).Create(&u).Error; err != nil {
		return fmt.Errorf("创建默认管理员失败: %w", err)
	}
	log.Printf("已创建默认管理员账户: %s", cfg.Account)
	return nil
}

func principalOf(u *models.User) *models.Principal {
	return &models.Principal{
		ID:          u.ID,
		UserAccount: u.UserAccount,
		UserName:    u.UserName,
		UserRole:    u.UserRole,
	}
}
