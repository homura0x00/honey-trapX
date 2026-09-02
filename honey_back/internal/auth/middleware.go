package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"honey_back/internal/models"
	"honey_back/internal/utils/res"
)

// LoginRequired 登录校验：从 Cookie 取 session id → 查 redis → 把登录态写入 ctx。
func (s *Service) LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie(s.session.CookieName)
		if err != nil || sid == "" {
			res.FailWithStatus(c, http.StatusUnauthorized, res.UserNotLogin, "未登录")
			return
		}

		p, err := s.CurrentSession(c.Request.Context(), sid)
		if err != nil {
			res.FailWithStatus(c, http.StatusUnauthorized, res.CodeOf(err), res.Message(err))
			return
		}

		c.Set(models.CtxPrincipal, *p)
		c.Next()
	}
}

// AdminRequired 管理员校验（需先经过 LoginRequired）。
func (s *Service) AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get(models.CtxPrincipal)
		if !ok {
			res.FailWithStatus(c, http.StatusUnauthorized, res.UserNotLogin, "未登录")
			return
		}
		p, ok := v.(models.Principal)
		if !ok || p.UserRole != models.RoleAdmin {
			res.FailWithStatus(c, http.StatusForbidden, res.Permission, "无权限")
			return
		}
		c.Next()
	}
}
