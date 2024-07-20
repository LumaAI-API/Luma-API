package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"luma-api/common"
	"luma-api/middleware"
	"net/http"
	"runtime/debug"
)

func main() {
	common.SetupLogger()
	common.Logger.Info("Luma-API started")

	if common.DebugEnabled {
		common.Logger.Info("running in debug mode")
		gin.SetMode(gin.ReleaseMode)
	}
	if common.PProfEnabled {
		common.SafeGoroutine(func() {
			log.Println(http.ListenAndServe("0.0.0.0:8005", nil))
		})
		common.Logger.Info("running in pprof")
	}
	common.InitTemplate()

	common.SafeGoroutine(func() {
		StartAllKeepAlive()
	})

	// Initialize HTTP server
	server := gin.New()
	server.Use(middleware.RequestId())
	server.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		common.Logger.Error(fmt.Sprintf("panic detected: %v, stack: %s", err, string(debug.Stack())))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": fmt.Sprintf("Panic detected, error: %v. Please contact site admin", err),
				"type":    "api_panic",
			},
		})
	}))
	server.Use(middleware.GinzapWithConfig())

	RegisterRouter(server)

	common.Logger.Info("Start: 0.0.0.0:" + common.Port)
	err := server.Run(":" + common.Port)
	if err != nil {
		common.Logger.Fatal("failed to start HTTP server: " + err.Error())
	}
}
