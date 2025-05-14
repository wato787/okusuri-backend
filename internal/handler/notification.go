package handler

import (
	"net/http"
	"okusuri-backend/internal/dto"
	"okusuri-backend/internal/model"
	"okusuri-backend/internal/repository"
	"okusuri-backend/pkg/helper"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationRepo *repository.NotificationRepository
	userRepo         *repository.UserRepository
}

func NewNotificationHandler(notificationRepo *repository.NotificationRepository, userRepo *repository.UserRepository) *NotificationHandler {
	return &NotificationHandler{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

// GetSetting は通知設定を取得するハンドラー
func (h *NotificationHandler) GetSetting(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// 通知設定を取得
	setting, err := h.notificationRepo.GetSettingByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notification setting"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// RegisterSetting は通知設定を登録するハンドラー
func (h *NotificationHandler) RegisterSetting(c *gin.Context) {
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
	setting := model.NotificationSetting{
		UserID:    userID,
		IsEnabled: req.IsEnabled,
		FcmToken:  req.FcmToken,
		Platform:  req.Platform,
	}

	// リポジトリに登録処理を依頼
	if err := h.notificationRepo.RegisterSetting(&setting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register notification setting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification setting registered successfully"})
}

// SendNotification は通知を送信するハンドラー
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	// ユーザーを全取得
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}

	// 通知設定を全取得
	settings, err := h.notificationRepo.GetAllSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notification settings"})
		return
	}

	// 通知設定をユーザーに紐づける
	settingsMap := make(map[string]model.NotificationSetting)
	for _, setting := range settings {
		settingsMap[setting.UserID] = setting
	}

	// リポジトリに通知送信処理を依頼
	for _, user := range users {
		setting, ok := settingsMap[user.ID]
		if !ok {
			continue
		}
		if setting.IsEnabled {
			err := h.notificationRepo.SendNotification(user, setting, "お薬の時間です🐣")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send notification"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification sent successfully"})
}
