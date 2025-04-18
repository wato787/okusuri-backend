package routes

import (
	"okusuri-backend/controller"

	"github.com/gin-gonic/gin"
)

func RegisterNotificationRoutes(router *gin.RouterGroup, notificationController *controller.NotificationController) {
	notification := router.Group("/notification")
	{
		// ユーザーのFCM通知設定を取得
		notification.GET("/setting", notificationController.GetNotificationSetting)
		// ユーザーのFCM通知設定を登録
		notification.POST("/setting", notificationController.RegisterNotificationSetting)
	}
}
