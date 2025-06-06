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
	medicationRepo   *repository.MedicationRepository
	medicationSvc    *service.MedicationService
}

func NewNotificationHandler(
	notificationRepo *repository.NotificationRepository,
	userRepo *repository.UserRepository,
	notificationSvc *service.NotificationService,
	medicationRepo *repository.MedicationRepository,
	medicationSvc *service.MedicationService,
) *NotificationHandler {
	return &NotificationHandler{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		notificationSvc:  notificationSvc,
		medicationRepo:   medicationRepo,
		medicationSvc:    medicationSvc,
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
	requestTime := time.Now()
	h.logRequestStart(c, requestTime)

	users, settings, err := h.fetchUsersAndSettings(c)
	if err != nil {
		return
	}

	settingsMap := h.buildSettingsMap(settings)
	sentCount := h.processNotifications(users, settingsMap)

	h.logAndRespond(c, requestTime, sentCount)
}

// logRequestStart はリクエスト開始時のログを出力する
func (h *NotificationHandler) logRequestStart(c *gin.Context, requestTime time.Time) {
	fmt.Printf("\n========== 通知送信処理開始 [%s] ==========\n", requestTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("リクエストパス: %s\n", c.Request.URL.Path)
	fmt.Printf("リクエスト元IP: %s\n", c.ClientIP())
	fmt.Printf("リクエストID: %s\n", c.Writer.Header().Get("Request-ID"))
}

// fetchUsersAndSettings はユーザーと通知設定を取得する
func (h *NotificationHandler) fetchUsersAndSettings(c *gin.Context) ([]model.User, []model.NotificationSetting, error) {
	users, userErr := h.userRepo.GetAllUsers()
	if userErr != nil {
		fmt.Printf("エラー: ユーザー取得失敗: %v\n", userErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return nil, nil, userErr
	}
	fmt.Printf("取得したユーザー数: %d\n", len(users))

	settings, settingsErr := h.notificationRepo.GetAllSettings()
	if settingsErr != nil {
		fmt.Printf("エラー: 通知設定取得失敗: %v\n", settingsErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notification settings"})
		return nil, nil, settingsErr
	}
	fmt.Printf("取得した通知設定数: %d\n", len(settings))

	return users, settings, nil
}

// buildSettingsMap は通知設定をユーザーIDでマップ化する
func (h *NotificationHandler) buildSettingsMap(
	settings []model.NotificationSetting,
) map[string]model.NotificationSetting {
	settingsMap := make(map[string]model.NotificationSetting)
	for _, setting := range settings {
		existingSetting, exists := settingsMap[setting.UserID]
		if !exists || setting.UpdatedAt.After(existingSetting.UpdatedAt) {
			settingsMap[setting.UserID] = setting
		}
	}
	fmt.Printf("通知対象ユーザー数: %d\n", len(settingsMap))
	return settingsMap
}

// processNotifications は各ユーザーに通知を送信する
func (h *NotificationHandler) processNotifications(
	users []model.User, settingsMap map[string]model.NotificationSetting,
) int {
	sentSubs := make(map[string]bool)
	fmt.Println("----- 通知送信処理開始 -----")

	for _, user := range users {
		if h.sendUserNotification(user, settingsMap, sentSubs) {
			fmt.Printf("ユーザーID: %s への通知送信成功\n", user.ID)
		}
	}

	fmt.Printf("----- 通知送信処理完了: 合計%d件送信 -----\n", len(sentSubs))
	return len(sentSubs)
}

// sendUserNotification は個別ユーザーに通知を送信する
func (h *NotificationHandler) sendUserNotification(
	user model.User, settingsMap map[string]model.NotificationSetting, sentSubs map[string]bool,
) bool {
	setting, ok := settingsMap[user.ID]
	if !ok || !setting.IsEnabled {
		return false
	}

	if _, alreadySent := sentSubs[setting.Subscription]; alreadySent && setting.Subscription != "" {
		return false
	}

	message := h.getNotificationMessage(user.ID)
	medicationStatus, statusErr := h.medicationSvc.GetMedicationStatus(user.ID)
	if statusErr == nil {
		message = h.generateStatusBasedMessage(medicationStatus)
	}

	consecutiveDays := 0
	if statusErr == nil {
		consecutiveDays = medicationStatus.CurrentStreak
	}

	sendErr := h.notificationSvc.SendNotificationWithDays(user, setting, message, consecutiveDays)
	if sendErr != nil {
		fmt.Printf("エラー: 通知送信失敗: %v\n", sendErr)
		return false
	}

	if setting.Subscription != "" {
		sentSubs[setting.Subscription] = true
	}
	return true
}

// getNotificationMessage はデフォルトの通知メッセージを取得する
func (h *NotificationHandler) getNotificationMessage(userID string) string {
	return "お薬の時間です。忘れずに服用してください。"
}

// logAndRespond は処理結果をログ出力してレスポンスを返す
func (h *NotificationHandler) logAndRespond(c *gin.Context, requestTime time.Time, sentCount int) {
	processingTime := time.Since(requestTime)
	fmt.Printf("処理時間: %v\n", processingTime)

	c.JSON(http.StatusOK, gin.H{
		"message":         "notification sent successfully",
		"sent_count":      sentCount,
		"process_time_ms": processingTime.Milliseconds(),
	})
	fmt.Printf("========== 通知送信処理終了 [%s] ==========\n\n",
		time.Now().Format("2006-01-02 15:04:05"))
}

// generateStatusBasedMessage はユーザーの薬のステータスに応じた通知メッセージを生成する
func (h *NotificationHandler) generateStatusBasedMessage(status *dto.MedicationStatusResponse) string {
	if status.IsRestPeriod {
		// 休薬期間中のメッセージ
		if status.RestDaysLeft > 0 {
			return fmt.Sprintf("現在休薬期間中です。あと%d日で服薬を再開してください。", status.RestDaysLeft)
		} else {
			return "休薬期間が終了しました。本日から服薬を再開してください。"
		}
	} else {
		// 通常の服薬期間のメッセージ
		if status.CurrentStreak > 0 {
			return fmt.Sprintf("お薬の時間です。忘れずに服用してください。（連続%d日目）", status.CurrentStreak)
		} else {
			return "お薬の時間です。忘れずに服用してください。"
		}
	}
}
