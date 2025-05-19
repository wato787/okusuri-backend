package dto

// 通知設定リクエスト
type RegisterNotificationSettingRequest struct {
	Subscription string `json:"subscription" binding:"required"` // FCMTokenをSubscriptionに変更
	IsEnabled    bool   `json:"isEnabled" binding:"required"`
	Platform     string `json:"platform" binding:"required"`
}
