package middleware

import (
	"fmt"
	"github.com/985492783/sparrow-go/pkg/config"
	"github.com/gin-gonic/gin"
)

type AuthMiddleWare struct {
	config *config.SparrowConfig
	permit string
}

func Auth(permit string, config *config.SparrowConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.AuthEnabled {
			var user struct {
				username string `form:"username"`
				password string `form:"password"`
			}
			err := c.Bind(&user)
			if err != nil {
				c.Writer.WriteHeader(301)
				c.Writer.WriteString(fmt.Sprintf("server unavailable %s", err))
				return
			}
			err = config.Authority(user.username, user.password, permit)
			if err != nil {
				c.Writer.WriteHeader(401)
				c.Writer.WriteString(fmt.Sprintf("403 Unauthorized %s", err))
				return
			}
		}

		c.Next()

	}
}

type ConsoleResponse gin.H

func NewResponseSuccess(data any) ConsoleResponse {
	return ConsoleResponse{
		"code": 200,
		"body": data,
	}
}

func NewResponseFailed(code int, err error) *ConsoleResponse {
	return &ConsoleResponse{
		"code":   code,
		"errMsg": err.Error(),
	}
}
