package dto

type RegisterNotificationSettingRequest struct {
	FcmToken  string `json:"fcm_token" binding:"required"`
	IsEnabled bool   `json:"is_enabled" binding:"required"`
	Platform  string `json:"platform" binding:"required"`
}
