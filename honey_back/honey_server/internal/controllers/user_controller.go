package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"honey_back/honey_server/config"
	"honey_back/honey_server/global"
	res2 "honey_back/honey_server/internal/common/res"
	"honey_back/honey_server/internal/models"
	"honey_back/honey_server/internal/models/dto"
	"honey_back/honey_server/internal/models/vo"
	"honey_back/honey_server/internal/query"
	"honey_back/honey_server/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUser(ctx *gin.Context) {
	var user dto.UserRegisterReq
	// 1. 注册信息校验
	if err := ctx.ShouldBindJSON(&user); err != nil {
		fmt.Println("参数绑定错误：", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误：" + err.Error()})
	}
	if len(user.UserAccount) < 6 || len(user.UserPassword) > 20 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "账户名必须为 6-20 个字符"})
	}
	if len(user.UserPassword) < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "密码必须至少 8 个字符"})
	}
	if user.UserPassword != user.CheckPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "两次输入的密码不一致"})
	}

	// 2. query user
	u := query.Use(global.MysqlDb)
	existingUser, err := u.User.Where(u.User.UserAccount.Eq(user.UserAccount)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("账户查询错误：", err.Error())
		res2.Error(ctx, res2.SystemError, "查询账户失败，请重试")
		return
	}
	// 若查到记录（existingUser 非 nil），说明账户已存在
	if existingUser != nil {
		res2.Error(ctx, res2.ParamCode, "该账户已被注册，请更换账户名")
		return
	}

	// 3. 密码加密
	hash, err := utils.HashPassword(user.UserPassword)
	user.UserPassword = hash // 暂存
	if err != nil {
		res2.Error(ctx, res2.SystemError, "密码未通过")
		return
	}

	// 4.存库
	newUser := models.User{
		UserAccount:  user.UserAccount,
		UserPassword: user.UserPassword,
		UserName:     user.UserName,
		UserRole:     "USER",
	}
	// 执行存库操作
	if err := u.User.Create(&newUser); err != nil {
		fmt.Println("用户存库错误：", err.Error())
		res2.Error(ctx, res2.SystemError, "注册失败，请重试")
		return
	}
	//jwtToken, err := utils.GenerateJWT(newUser.UserAccount)
	//if err != nil {
	//	res.Error(ctx, res.SystemError, "Token 生成失败："+err.Error())
	//	return
	//}
	// 5. 格式化返回脱敏后的新用户信息
	ctx.JSON(http.StatusOK, gin.H{
		"success": newUser,
	})
}

// LoginUser 用户登陆
func LoginUser(c *gin.Context) {
	var loginUser dto.LoginUserReq
	// 1. 登陆信息校验
	if err := c.ShouldBind(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		fmt.Println("参数绑定错误：", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误：" + err.Error()})
	}
	if len(loginUser.UserAccount) < 6 || len(loginUser.UserPassword) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "账户名必须为 6-20 个字符"})
	}
	if len(loginUser.UserPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码必须至少 8 个字符"})
	}

	// 2. 查库
	u := query.Use(global.MysqlDb)
	existingUser, err := u.User.Where(u.User.UserAccount).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("账户查询错误：", err.Error())
		res2.Error(c, res2.SystemError, "查询账户失败，请重试")
		return
	}

	// 返回给前端的用户（脱敏后）数据
	user := vo.UserVO{
		Id:          existingUser.ID,
		UserAccount: existingUser.UserAccount,
		Username:    existingUser.UserName,
		UserRole:    existingUser.UserRole,
	}

	// 随机生成session_id
	sessionId := utils.GenerateSession()
	userJson, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	// 3. 存 redis
	redisKey := config.AppConfig.Session.SessionRedisPrefix + sessionId
	// 将 int 类型的数值 -> time 类型 -> 小时制单位
	sessionExpire := time.Duration(config.AppConfig.Session.SessionExpire) * time.Hour
	// session 存 redis
	ctx := c.Request.Context()
	err = global.RedisDb.Set(ctx, redisKey, userJson, sessionExpire).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	// 4. 返回给前端的 Cookie
	c.SetCookie(
		config.AppConfig.Session.SessionCookieKey, // Cookie 键名
		sessionId,                                 // Cookie值（session_id）
		int(sessionExpire.Seconds()),              // 过期时间（秒）
		"/",                                       // 路径（全站有效）
		"",                                        // 域名（单服务器置空，默认当前域名）
		false,                                     // Secure：仅HTTPS传输（生产环境必开）
		true,                                      // HttpOnly：禁止JS访问（防XSS）
	)
	// （可选）SameSite=Strict 防CSRF
	c.Header("Set-Cookie", c.Writer.Header().Get("Set-Cookie")+"; SameSite=Strict")
}

// GetLoginUser 获取当前登陆用户信息（脱敏后）
func GetLoginUser(c *gin.Context) {
	// 从上下文获取用户信息
	user, exists := c.Get("currentUser")
	if !exists {
		// 不存在：用户未登录或会话过期，返回401
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或会话已过期"})
		return
	}
	// 安全类型断言：避免类型不匹配导致的panic
	currentUser, ok := user.(vo.UserVO)
	if !ok {
		// 类型不匹配：中间件存入的数据类型错误，属于服务器内部问题
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息格式错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": currentUser})
}

// UserLogout 用户注销（退出登陆）
func UserLogout(c *gin.Context) {
	// 1. Get Cookie
	sessionID, err := c.Cookie(config.AppConfig.Session.SessionCookieKey)
	if err != nil || sessionID == "" {
		c.JSON(http.StatusOK, gin.H{"message": "已登出"})
		return
	}

	// 2. Delete session-redis
	redisKey := config.AppConfig.Session.SessionRedisPrefix + sessionID
	global.RedisDb.Del(c, redisKey)

	// 3. Delete Cookie
	c.SetCookie("token", sessionID, -1, "/", "", false, true)
}

// ------------- 管理用户（管理员） --------------

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	var user dto.UpdateUserReq
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数不能为空"})
	}
	q := query.Use(global.MysqlDb)
	_, err := q.User.Where(q.User.ID.Eq(user.UserId)).First()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
	}

	updates, err := q.User.Where(q.User.ID.Eq(user.UserId)).Updates(models.User{
		UserName:     user.UserName,
		UserPassword: user.UserPassword,
		UserRole:     user.UserRole,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
	}

	c.JSON(http.StatusOK, gin.H{"success": updates})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	var userId dto.Req
	// 1. 参数绑定
	if err := c.BindJSON(&userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不能为空"})
		return
	}

	// 2. 查库
	q := query.Use(global.MysqlDb)
	user, err := q.User.Where(q.User.ID.Eq(userId.Id)).First()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
	} else if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
	}

	// 3. 删除
	res, err := q.User.Delete(&models.User{ID: userId.Id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"success": res})
}
