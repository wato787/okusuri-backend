package routes

import (
	"okusuri-backend/controller"
	"okusuri-backend/internal/middleware"
	"okusuri-backend/repository"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	// controllerを初期化
	notificationController := controller.NewNotificationController()
	medicationLogController := controller.NewMedicationLogController()

	// userRepositoryを初期化
	userRepository := repository.NewUserRepository()

	// Ginのルーターを作成
	router := gin.Default()

	// グローバルミドルウェアの設定
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})

		api.POST(("/notification"), notificationController.SendNotification)

		notificationSetting := api.Group("/notification/setting")
		notificationSetting.Use(middleware.Auth(userRepository))
		notificationSetting.GET("/", notificationController.GetNotificationSetting)
		notificationSetting.POST("/", notificationController.RegisterNotificationSetting)

		medicationLog := api.Group("/medication-log")
		medicationLog.Use(middleware.Auth(userRepository))
		medicationLog.POST("/", medicationLogController.RegisterMedicationLog)
		medicationLog.GET("/", medicationLogController.GetMedicationLogs)
	}

	return router
}
