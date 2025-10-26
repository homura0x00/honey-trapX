package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HostFilterConfig struct {
	Hosts         []string
	IsWhitelisted bool
}

func NewHostFilter(config HostFilterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前请求路由 Host (自动去除端口，只保留)
		currentHost := c.Request.Host

		inList := false
		for _, host := range config.Hosts {
			if currentHost == host {
				inList = true
				break
			}
		}

		if config.IsWhitelisted {
			if !inList {
				// 白名单模式，不在列表则拦截
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Host not allowed.",
					"hosts": currentHost,
				})
			} else {
				// 黑名单模式：拦截列表里的Host
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Host blocked！",
					"hosts": currentHost,
				})
			}
		}

		c.Next()
	}
}

func SecurityHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	}
}
