package middleware

import (
	"encoding/json"
	"errors"
	"honey_back/honey_server/config"
	"honey_back/honey_server/global"
	"honey_back/honey_server/internal/models/enums"
	"honey_back/honey_server/internal/models/vo"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// AdminMiddleware 鉴权中间件（鉴别是否为ADMIN）
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
			return
		}
		currentUser := user.(vo.UserVO)
		if currentUser.UserRole != enums.UserRoleAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "非特权用户"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AuthMiddleware 全局登陆验证
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Cookie 中获取 session_id
		sessionID, err := c.Cookie(config.AppConfig.Session.SessionCookieKey)
		if err != nil || sessionID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session cookie"})
			c.Abort()
			return
		}

		// 2. 查 Redis
		sessionCookiePrefix := config.AppConfig.Session.SessionRedisPrefix
		redisKey := sessionCookiePrefix + sessionID
		ctx := c.Request.Context()
		userJSON, err := global.RedisDb.Get(ctx, redisKey).Result()
		if errors.Is(err, redis.Nil) {
			// session_id 不存在 or 过期
			c.JSON(http.StatusUnauthorized, gin.H{"error": "登陆已过期，请重新登陆"})
			c.Abort()
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
			c.Abort()
			return
		}

		// 3. 解析用户信息，供后续接口使用
		var user vo.UserVO
		if err = json.Unmarshal([]byte(userJSON), &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		}
		// 存入上下文
		c.Set("currentUser", user)

		// （可选）4. 每次访问刷新 Redis-Cookie-Expire 时间
		// TODO 思考：是否要刷新Cookie时间
		global.RedisDb.Expire(ctx, redisKey, time.Duration(config.AppConfig.Session.SessionExpire)*time.Second)

		// 执行下一步操作
		c.Next()
	}
}

// GlobalMiddleware 全局操作记录器
func GlobalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO
		err := c.Err()
		if err != nil {
			log.Fatalln(err.Error())
		}
		c.Next()
	}
}
