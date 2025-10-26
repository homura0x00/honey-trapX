package res

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func response(c *gin.Context, status int, r Response) {
	c.JSON(status, r)
}

func Success(c *gin.Context, data interface{}) {
	response(c, http.StatusOK, Response{
		Code:    0,
		Data:    data,
		Message: "",
	})
}

func Error(c *gin.Context, status int, err string) {
	response(c, status, Response{
		Message: err,
		Data:    nil,
	})
}
