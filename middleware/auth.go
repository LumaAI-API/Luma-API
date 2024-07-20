package middleware

import (
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"detail": map[string]any{
					"reason": "unauthorized secret token",
					"code":   1,
				},
			})
			c.Abort()
			return
		}
	}
}
