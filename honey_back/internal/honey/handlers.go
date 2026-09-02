package honey

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"honey_back/internal/models"
	"honey_back/internal/utils/res"
)

type Handlers struct {
	svc *Service
}

func NewHandlers(svc *Service) *Handlers {
	return &Handlers{svc: svc}
}

func principalID(c *gin.Context) (int64, bool) {
	v, ok := c.Get(models.CtxPrincipal)
	if !ok {
		return 0, false
	}
	p, ok := v.(models.Principal)
	return p.ID, ok
}

// GET /api/images
func (h *Handlers) ListImages(c *gin.Context) {
	list, err := h.svc.ListImages(c.Request.Context())
	if err != nil {
		res.Fail(c, err)
		return
	}
	res.Ok(c, list)
}

// POST /api/deployments
func (h *Handlers) CreateDeployment(c *gin.Context) {
	uid, ok := principalID(c)
	if !ok {
		res.Fail(c, res.E(res.UserNotLogin, "未登录"))
		return
	}
	var req CreateDeploymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		res.Fail(c, res.E(res.ParamCode, "参数格式错误"))
		return
	}
	view, err := h.svc.Create(c.Request.Context(), uid, &req)
	if err != nil {
		res.Fail(c, err)
		return
	}
	res.Ok(c, view)
}

// GET /api/deployments
func (h *Handlers) ListDeployments(c *gin.Context) {
	uid, ok := principalID(c)
	if !ok {
		res.Fail(c, res.E(res.UserNotLogin, "未登录"))
		return
	}
	list, err := h.svc.List(c.Request.Context(), uid)
	if err != nil {
		res.Fail(c, err)
		return
	}
	res.Ok(c, list)
}

// GET /api/deployments/:id
func (h *Handlers) GetDeployment(c *gin.Context) {
	uid, ok := principalID(c)
	if !ok {
		res.Fail(c, res.E(res.UserNotLogin, "未登录"))
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		res.Fail(c, res.E(res.ParamCode, "参数格式错误"))
		return
	}
	view, err := h.svc.Get(c.Request.Context(), uid, id)
	if err != nil {
		res.Fail(c, err)
		return
	}
	res.Ok(c, view)
}

// POST /api/deployments/:id/start
func (h *Handlers) StartDeployment(c *gin.Context) {
	uid, ok := principalID(c)
	if !ok {
		res.Fail(c, res.E(res.UserNotLogin, "未登录"))
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		res.Fail(c, res.E(res.ParamCode, "参数格式错误"))
		return
	}
	if err := h.svc.Start(c.Request.Context(), uid, id); err != nil {
		res.Fail(c, err)
		return
	}
	res.OkMsg(c, nil, "已启动")
}

// POST /api/deployments/:id/stop
func (h *Handlers) StopDeployment(c *gin.Context) {
	uid, ok := principalID(c)
	if !ok {
		res.Fail(c, res.E(res.UserNotLogin, "未登录"))
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		res.Fail(c, res.E(res.ParamCode, "参数格式错误"))
		return
	}
	if err := h.svc.Stop(c.Request.Context(), uid, id); err != nil {
		res.Fail(c, err)
		return
	}
	res.OkMsg(c, nil, "已停止")
}

// DELETE /api/deployments/:id
func (h *Handlers) DeleteDeployment(c *gin.Context) {
	uid, ok := principalID(c)
	if !ok {
		res.Fail(c, res.E(res.UserNotLogin, "未登录"))
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		res.Fail(c, res.E(res.ParamCode, "参数格式错误"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uid, id); err != nil {
		res.Fail(c, err)
		return
	}
	res.OkMsg(c, nil, "已删除")
}
