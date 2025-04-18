package model

import (
	"gorm.io/gorm"
)

// ユーザーのFCM通知設定を管理する構造体
type NotificationSetting struct {
	gorm.Model
	UserID    string `json:"user_id" gorm:"not null;uniqueIndex:idx_user_id"`
	IsEnabled bool   `json:"is_enabled" gorm:"default:true"`
	FcmToken  string `json:"fcm_token" gorm:"not null;uniqueIndex:idx_fcm_token"`
	Platform  string `json:"platform" gorm:"not null;uniqueIndex:idx_platform"`
}
