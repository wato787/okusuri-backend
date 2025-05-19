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
		UserID:       userID,
		IsEnabled:    req.IsEnabled,
		Platform:     req.Platform,
		Subscription: req.Subscription,
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
		subPreview := getPreview(setting.Subscription)
		fmt.Printf("  %d. ユーザーID: %s, 有効: %t, サブスクリプション: %s, 更新日時: %s\n",
			i+1, setting.UserID, setting.IsEnabled, subPreview, setting.UpdatedAt)
	}

	// 通知設定をユーザーに紐づける - 各ユーザーの最新設定のみを保持
	settingsMap := make(map[string]model.NotificationSetting)
	for _, setting := range settings {
		existingSetting, exists := settingsMap[setting.UserID]
		if !exists || setting.UpdatedAt.After(existingSetting.UpdatedAt) {
			settingsMap[setting.UserID] = setting
			fmt.Printf("ユーザーID: %s の通知設定を登録/更新 (サブスクリプション: %s)\n",
				setting.UserID, getPreview(setting.Subscription))
		}
	}
	fmt.Printf("通知対象ユーザー数: %d\n", len(settingsMap))

	// 一時的に送信済みのサブスクリプションを記録するセット
	sentSubs := make(map[string]bool)

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

		// 送信済みのサブスクリプションをスキップ
		if _, alreadySent := sentSubs[setting.Subscription]; alreadySent && setting.Subscription != "" {
			fmt.Printf("ユーザーID: %s のサブスクリプションはすでに送信済みのためスキップします (サブスクリプション: %s)\n",
				user.ID, getPreview(setting.Subscription))
			continue
		}

		// 通知送信
		fmt.Printf("ユーザーID: %s に通知送信中 (サブスクリプション: %s)\n",
			user.ID, getPreview(setting.Subscription))

		// 服薬メッセージは簡単なものにしておく
		message := "お薬の時間です。忘れずに服用してください。"

		err := h.notificationSvc.SendNotification(user, setting, message)
		if err != nil {
			fmt.Printf("エラー: 通知送信失敗: %v\n", err)
			// エラーがあっても処理を続行
			continue
		}

		// 送信済みとしてマーク (空のサブスクリプションはマークしない)
		if setting.Subscription != "" {
			sentSubs[setting.Subscription] = true
		}
		fmt.Printf("ユーザーID: %s への通知送信成功\n", user.ID)
	}
	fmt.Printf("----- 通知送信処理完了: 合計%d件送信 -----\n", len(sentSubs))

	// 処理時間を計算
	processingTime := time.Since(requestTime)
	fmt.Printf("処理時間: %v\n", processingTime)

	c.JSON(http.StatusOK, gin.H{
		"message":         "notification sent successfully",
		"sent_count":      len(sentSubs),
		"process_time_ms": processingTime.Milliseconds(),
	})
	fmt.Printf("========== 通知送信処理終了 [%s] ==========\n\n",
		time.Now().Format("2006-01-02 15:04:05"))
}

// トークンやサブスクリプションの先頭数文字を取得するヘルパー関数
func getPreview(str string) string {
	if str == "" {
		return "空"
	}

	if len(str) <= 10 {
		return str
	}
	return str[:10] + "..."
}
