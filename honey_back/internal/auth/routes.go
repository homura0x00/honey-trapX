package auth

import "github.com/gin-gonic/gin"

// Register 注册 /api/auth 路由
func Register(api *gin.RouterGroup, h *Handlers, svc *Service) {
	api.POST("/auth/login", h.Login)

	authed := api.Group("/auth", svc.LoginRequired())
	authed.POST("/logout", h.Logout)
	authed.GET("/me", h.Me)
}
