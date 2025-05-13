package notification

import (
	"net/http"
	"okusuri-backend/internal/common/user"
	"okusuri-backend/pkg/helper"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service     *Service
	userService *user.Service
}

// NewHandler は通知ハンドラーの新しいインスタンスを作成する
func NewHandler() *Handler {
	repository := NewRepository()
	service := NewService(repository)
	userRepository := user.NewRepository()
	userService := user.NewService(userRepository)

	return &Handler{
		service:     service,
		userService: userService,
	}
}

// GetSetting は通知設定を取得するハンドラー
func (h *Handler) GetSetting(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// 通知設定を取得
	setting, err := h.service.GetSettingByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notification setting"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// RegisterSetting は通知設定を登録するハンドラー
func (h *Handler) RegisterSetting(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// リクエストボディから通知設定を取得
	var req RegisterNotificationSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// 通知設定をモデルに変換
	setting := NotificationSetting{
		UserID:    userID,
		IsEnabled: req.IsEnabled,
		FcmToken:  req.FcmToken,
		Platform:  req.Platform,
	}

	// サービス層に登録処理を依頼
	if err := h.service.RegisterSetting(&setting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register notification setting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification setting registered successfully"})
}

// SendNotification は通知を送信するハンドラー
func (h *Handler) SendNotification(c *gin.Context) {
	// ユーザーを全取得
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}

	// 通知設定を全取得
	settings, err := h.service.GetAllSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notification settings"})
		return
	}

	// 通知設定をユーザーに紐づける
	settingsMap := make(map[string]NotificationSetting)
	for _, setting := range settings {
		settingsMap[setting.UserID] = setting
	}

	// サービス層に通知送信処理を依頼
	for _, user := range users {
		setting, ok := settingsMap[user.ID]
		if !ok {
			continue
		}
		if setting.IsEnabled {
			err := h.service.SendNotification(user, setting, "お薬の時間です🐣")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send notification"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification sent successfully"})
}
