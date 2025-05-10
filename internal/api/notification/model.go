package notification

import (
	"time"

	"gorm.io/gorm"
)

// ユーザーのFCM通知設定を管理する構造体
type NotificationSetting struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"` // gorm.DeletedAtに変更
	UserID    string         `json:"userId" gorm:"not null;index:idx_user_platform,unique:true,part:1"`
	Platform  string         `json:"platform" gorm:"not null;index:idx_user_platform,unique:true,part:2"`
	IsEnabled bool           `json:"isEnabled" gorm:"default:true"`
	FcmToken  string         `json:"fcmToken" gorm:"not null"`
}
