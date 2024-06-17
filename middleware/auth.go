package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"luma-api/common"
	"net/http"
	"strings"
)

func SecretAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		if common.SecretToken == "" {
			return
		}
		accessToken := c.Request.Header.Get("Authorization")
		accessToken = strings.TrimLeft(accessToken, "Bearer ")
		if accessToken == common.SecretToken {
			c.Next()
		} else {
			common.WrapperLumaError(c, fmt.Errorf("unauthorized secret token"), http.StatusUnauthorized)
			c.Abort()
			return
		}
	}
}
