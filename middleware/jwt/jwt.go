package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/luciferCN22/go-gin-example/pkg/e"
	"github.com/luciferCN22/go-gin-example/pkg/util"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			_, err := util.ParseToken(token)
			if errors.Is(err, jwt.ErrTokenExpired) {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			} else if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
