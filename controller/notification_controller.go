package controller

import (
	"net/http"
	"okusuri-backend/dto"
	"okusuri-backend/pkg/helper"

	"okusuri-backend/model"
	"okusuri-backend/repository"
	"okusuri-backend/service"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	NotificationRepository *repository.NotificationRepository
	NotificationService    *service.NotificationService
	UserRepository         *repository.UserRepository
	UserService            *service.UserService
}

// NewNotificationController は新しいNotificationControllerのインスタンスを作成する
func NewNotificationController() *NotificationController {

	notificationRepository := repository.NewNotificationRepository()
	notificationService := service.NewNotificationService(notificationRepository)
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository)
	return &NotificationController{
		NotificationRepository: notificationRepository,
		NotificationService:    notificationService,
		UserRepository:         userRepository,
		UserService:            userService,
	}
}

func (nc *NotificationController) GetNotificationSetting(c *gin.Context) {
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

func (nc *NotificationController) RegisterNotificationSetting(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// リクエストボディから通知設定を取得
	var req dto.RegisterNotificationSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// 通知設定をモデルに変換
	notificationSetting := model.NotificationSetting{
		UserID:    userID,
		IsEnabled: req.IsEnabled,
		FcmToken:  req.FcmToken,
		Platform:  req.Platform,
	}

	// サービス層に登録処理を依頼
	if err := nc.NotificationService.RegisterNotificationSetting(&notificationSetting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register notification setting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification setting registered successfully"})
}

// 通知送信
func (nc *NotificationController) SendNotification(c *gin.Context) {
	// ユーザーを全取得
	users, err := nc.UserService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}

	// 通知設定を全取得
	notificationSettings, err := nc.NotificationService.GetAllNotificationSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notification settings"})
		return
	}
	// 通知設定をユーザーに紐づける
	notificationSettingsMap := make(map[string]model.NotificationSetting)
	for _, ns := range notificationSettings {
		notificationSettingsMap[ns.UserID] = ns
	}
	// サービス層に通知送信処理を依頼
	for _, user := range users {
		notificationSetting, ok := notificationSettingsMap[user.ID]
		if !ok {
			continue
		}
		if notificationSetting.IsEnabled {
			err := nc.NotificationService.SendNotification(user, notificationSetting, "お薬の時間です")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send notification"})
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "notification sent successfully"})
}
