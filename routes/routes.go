package routes

import (
	"okusuri-backend/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	// コントローラを初期化
	authController := controller.NewAuthController()

	// Ginのルーターを作成
	router := gin.Default()
	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})
		RegisterAuthRoutes(api, authController)
	}

	return router
}
