package server

import (
	"github.com/gin-gonic/gin"

	"honey_back/internal/auth"
	"honey_back/internal/config"
	"honey_back/internal/honey"
)

// New 组装 gin 引擎：全局中间件 → /api 分组 → 各 feature 注册。
func New(cfg *config.Config, authSvc *auth.Service, honeySvc *honey.Service) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	if cfg.Server.FrontendURL != "" {
		r.Use(cors(cfg.Server.FrontendURL))
	}

	api := r.Group("/api")
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	auth.Register(api, auth.NewHandlers(authSvc), authSvc)
	honey.Register(api, honeySvc, authSvc.LoginRequired())

	return r
}

// cors 仅放行配置的前端来源，并允许携带 Cookie。
func cors(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Origin") == origin {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
