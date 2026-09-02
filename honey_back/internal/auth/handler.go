package auth

import (
	"log"
	"net/http"

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

// Login POST /api/auth/login
func (h *Handlers) Login(c *gin.Context) {
	var req struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		res.Fail(c, res.E(res.ParamCode, "参数格式错误"))
		return
	}

	p, err := h.svc.Login(c.Request.Context(), req.Account, req.Password)
	if err != nil {
		res.Fail(c, err)
		return
	}

	sid, err := h.svc.CreateSession(c.Request.Context(), p)
	if err != nil {
		log.Printf("create session: %v", err)
		res.Fail(c, res.E(res.SystemError, "服务器错误"))
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     h.svc.session.CookieName,
		Value:    sid,
		Path:     "/",
		MaxAge:   int(h.svc.ttl().Seconds()),
		HttpOnly: true,
		Secure:   h.svc.session.CookieSecure,
		SameSite: http.SameSiteStrictMode,
	})
	res.Ok(c, p)
}

// Logout POST /api/auth/logout
func (h *Handlers) Logout(c *gin.Context) {
	sid, _ := c.Cookie(h.svc.session.CookieName)
	h.svc.DestroySession(c.Request.Context(), sid)

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     h.svc.session.CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	res.OkMsg(c, nil, "已登出")
}

// Me GET /api/auth/me
func (h *Handlers) Me(c *gin.Context) {
	v, ok := c.Get(models.CtxPrincipal)
	if !ok {
		res.Fail(c, res.E(res.UserNotLogin, "未登录"))
		return
	}
	res.Ok(c, v)
}
