package dto

type RegisterNotificationSettingRequest struct {
	FcmToken  string `json:"fcmToken" binding:"required"`
	IsEnabled bool   `json:"isEnabled" binding:"required"`
	Platform  string `json:"platform" binding:"required"`
}
