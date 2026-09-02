package res

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Err 业务错误：handler 层用 CodeOf/Message 翻译，其余错误回通用 50000
type Err struct {
	Code int
	Msg  string
}

func (e *Err) Error() string { return e.Msg }

// E 构造业务错误
func E(code int, msg string) error { return &Err{Code: code, Msg: msg} }

// CodeOf 提取错误码：nil→成功，*Err→透传，其余→系统错误
func CodeOf(err error) int {
	if err == nil {
		return CodeOK
	}
	var e *Err
	if errors.As(err, &e) {
		return e.Code
	}
	return SystemError
}

// Message 提取用户可见消息：*Err→透传，其余→通用文案
func Message(err error) string {
	if err == nil {
		return ""
	}
	var e *Err
	if errors.As(err, &e) {
		return e.Msg
	}
	return "服务器错误"
}

func Ok(c *gin.Context, data interface{}) {
	c.JSON(200, Response{Code: CodeOK, Data: data})
}

func OkMsg(c *gin.Context, data interface{}, msg string) {
	c.JSON(200, Response{Code: CodeOK, Data: data, Message: msg})
}

// Fail 统一错误出口：HTTP 恒 200，业务错误靠 code 区分（沿用现有约定）
func Fail(c *gin.Context, err error) {
	c.JSON(200, Response{Code: CodeOf(err), Data: nil, Message: Message(err)})
}

// FailWithStatus 供中间件使用（如 401/403）
func FailWithStatus(c *gin.Context, status int, code int, msg string) {
	c.AbortWithStatusJSON(status, Response{Code: code, Data: nil, Message: msg})
}
