package handler

import (
	"fmt"
	"net/http"
	"okusuri-backend/internal/dto"
	"okusuri-backend/internal/model"
	"okusuri-backend/internal/repository"
	"okusuri-backend/internal/service"
	"okusuri-backend/pkg/helper"
	"time"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationRepo *repository.NotificationRepository
	userRepo         *repository.UserRepository
	notificationSvc  *service.NotificationService
}

func NewNotificationHandler(
	notificationRepo *repository.NotificationRepository,
	userRepo *repository.UserRepository,
	notificationSvc *service.NotificationService,
) *NotificationHandler {
	return &NotificationHandler{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		notificationSvc:  notificationSvc,
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
	// 詳細なログ出力
	requestTime := time.Now()
	fmt.Printf("\n========== 通知送信処理開始 [%s] ==========\n", requestTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("リクエストパス: %s\n", c.Request.URL.Path)
	fmt.Printf("リクエスト元IP: %s\n", c.ClientIP())
	fmt.Printf("リクエストID: %s\n", c.Writer.Header().Get("Request-ID"))

	// ユーザーを全取得
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		fmt.Printf("エラー: ユーザー取得失敗: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}
	fmt.Printf("取得したユーザー数: %d\n", len(users))
	
	// ユーザーIDを表示
	fmt.Println("ユーザーID一覧:")
	for i, user := range users {
		fmt.Printf("  %d. %s\n", i+1, user.ID)
	}

	// 通知設定を全取得
	settings, err := h.notificationRepo.GetAllSettings()
	if err != nil {
		fmt.Printf("エラー: 通知設定取得失敗: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notification settings"})
		return
	}
	fmt.Printf("取得した通知設定数: %d\n", len(settings))
	
	// 通知設定詳細を表示
	fmt.Println("通知設定一覧:")
	for i, setting := range settings {
		tokenPreview := getTokenPreview(setting.FcmToken)
		fmt.Printf("  %d. ユーザーID: %s, 有効: %t, トークン: %s, 更新日時: %s\n", 
			i+1, setting.UserID, setting.IsEnabled, tokenPreview, setting.UpdatedAt)
	}

	// 通知設定をユーザーに紐づける - 各ユーザーの最新設定のみを保持
	settingsMap := make(map[string]model.NotificationSetting)
	for _, setting := range settings {
		existingSetting, exists := settingsMap[setting.UserID]
		if !exists || setting.UpdatedAt.After(existingSetting.UpdatedAt) {
			settingsMap[setting.UserID] = setting
			fmt.Printf("ユーザーID: %s の通知設定を登録/更新 (トークン: %s)\n", 
				setting.UserID, getTokenPreview(setting.FcmToken))
		}
	}
	fmt.Printf("通知対象ユーザー数: %d\n", len(settingsMap))

	// 一時的に送信済みのFCMトークンを記録するセット
	sentTokens := make(map[string]bool)

	// 通知送信処理
	fmt.Println("----- 通知送信処理開始 -----")
	for _, user := range users {
		fmt.Printf("ユーザーID: %s の処理\n", user.ID)
		setting, ok := settingsMap[user.ID]
		if !ok {
			fmt.Printf("ユーザーID: %s の通知設定が見つかりません\n", user.ID)
			continue
		}
		if !setting.IsEnabled {
			fmt.Printf("ユーザーID: %s の通知は無効化されています\n", user.ID)
			continue
		}
		
		// 送信済みのトークンをスキップ
		if _, alreadySent := sentTokens[setting.FcmToken]; alreadySent {
			fmt.Printf("ユーザーID: %s のトークンはすでに送信済みのためスキップします (トークン: %s)\n", 
				user.ID, getTokenPreview(setting.FcmToken))
			continue
		}

		fmt.Printf("ユーザーID: %s に通知送信中 (トークン: %s)\n", 
			user.ID, getTokenPreview(setting.FcmToken))
		err := h.notificationSvc.SendNotification(user, setting, "お薬の時間です🐣")
		if err != nil {
			fmt.Printf("エラー: 通知送信失敗: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send notification"})
			return
		}

		// 送信済みとしてマーク
		sentTokens[setting.FcmToken] = true
		fmt.Printf("ユーザーID: %s への通知送信成功\n", user.ID)
	}
	fmt.Printf("----- 通知送信処理完了: 合計%d件送信 -----\n", len(sentTokens))

	// 処理時間を計算
	processingTime := time.Since(requestTime)
	fmt.Printf("処理時間: %v\n", processingTime)

	c.JSON(http.StatusOK, gin.H{
		"message":    "notification sent successfully",
		"sent_count": len(sentTokens),
		"process_time_ms": processingTime.Milliseconds(),
	})
	fmt.Printf("========== 通知送信処理終了 [%s] ==========\n\n", 
		time.Now().Format("2006-01-02 15:04:05"))
}

// FCMトークンの先頭数文字を取得するヘルパー関数
func getTokenPreview(token string) string {
	if len(token) <= 10 {
		return token
	}
	return token[:10] + "..."
}
