package controller

import (
	"net/http"
	"okusuri-backend/helper"
	"okusuri-backend/repository"
	"okusuri-backend/service"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	NotificationRepository *repository.NotificationRepository
	NotificationService    *service.NotificationService
}

// NewNotificationController は新しいNotificationControllerのインスタンスを作成する
func NewNotificationController() *NotificationController {

	notificationRepository := repository.NewNotificationRepository()
	notificationService := service.NewNotificationService(notificationRepository)
	return &NotificationController{
		NotificationRepository: notificationRepository,
		NotificationService:    notificationService,
	}
}

func (nc *NotificationController) RegisterNotificationSetting(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// 通知設定を取得
	notificationSetting, err := nc.NotificationService.GetNotificationSettingByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notification setting"})
		return
	}

	c.JSON(http.StatusOK, notificationSetting)
}
