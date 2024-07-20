package main

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"luma-api/docs"
	"luma-api/middleware"
)

func RegisterRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	apiRouter := r.Group("/luma", middleware.SecretAuth())
	{
		//old 废弃
		apiRouter.POST("/generations/", Generations)

		// current
		//apiRouter.GET("/generations/*action", Fetch)

		// submit
		apiRouter.POST("/generations", Generations)
		apiRouter.POST("/generations/:task_id/extend", ExtentGenerations)
		apiRouter.POST("/generations/file_upload", Upload)

		// get data
		apiRouter.GET("/generations/:task_id", FetchByID)
		apiRouter.GET("/generations/:task_id/download_video_url", GetDownloadUrl)

		apiRouter.GET("/subscription/usage", Usage)
		apiRouter.GET("/users/me", Me)
	}
}
