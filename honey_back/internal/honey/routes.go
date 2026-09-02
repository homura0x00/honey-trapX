package honey

import "github.com/gin-gonic/gin"

// Register 注册 /api/images、/api/deployments 路由（均需登录）
func Register(api *gin.RouterGroup, svc *Service, loginRequired gin.HandlerFunc) {
	h := NewHandlers(svc)

	images := api.Group("/images", loginRequired)
	images.GET("", h.ListImages)

	deployments := api.Group("/deployments", loginRequired)
	deployments.POST("", h.CreateDeployment)
	deployments.GET("", h.ListDeployments)
	deployments.GET("/:id", h.GetDeployment)
	deployments.POST("/:id/start", h.StartDeployment)
	deployments.POST("/:id/stop", h.StopDeployment)
	deployments.DELETE("/:id", h.DeleteDeployment)
}
