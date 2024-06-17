package common

import "github.com/gin-gonic/gin"

func WrapperLumaError(c *gin.Context, err error, statusCode int) {
	Logger.Errorw("wrapper luma error", "statusCode", statusCode, "err", err)
	c.JSON(statusCode, gin.H{
		"detail": map[string]any{
			"reason": err.Error(),
			"code":   1,
		},
	})
}

func GetLumaAuth() string {
	if AccessToken != "" {
		return "access_token=" + AccessToken
	} else if COOKIE != "" {
		return COOKIE
	}
	return ""
}
