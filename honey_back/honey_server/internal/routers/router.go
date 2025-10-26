package routers

import (
	"github.com/gin-gonic/gin"
)

var apiRouterGroup *gin.RouterGroup

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Security Group
	//hostFilter := middleware.HostFilterConfig{
	//	Hosts:         []string{""},
	//	IsWhitelisted: true,
	//}
	//router.Use(middleware.NewHostFilter(hostFilter))
	//router.Use(middleware.SecurityHandler())

	apiRouterGroup = router.Group("/api")
	//apiRouterGroup.Use(middleware.GlobalMiddleware())
	// 路由分组
	UserRouter()
	return router
}
