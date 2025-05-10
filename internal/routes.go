package routes

import (
	"okusuri-backend/internal/api/medication"
	"okusuri-backend/internal/api/notification"
	"okusuri-backend/internal/common/user"
	"okusuri-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	notificationHandler := notification.NewHandler()
	medicationLogHandler := medication.NewHandler()

	// userRepositoryを初期化
	userRepository := user.NewRepository()

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

		api.POST(("/notification"), notificationHandler.SendNotification)

		notificationSetting := api.Group("/notification/setting")
		notificationSetting.Use(middleware.Auth(userRepository))
		notificationSetting.GET("/", notificationHandler.GetSetting)
		notificationSetting.POST("/", notificationHandler.RegisterSetting)

		medicationLog := api.Group("/medication-log")
		medicationLog.Use(middleware.Auth(userRepository))
		medicationLog.POST("/", medicationLogHandler.RegisterLog)
		medicationLog.GET("/", medicationLogHandler.GetLogs)
	}

	return router
}
