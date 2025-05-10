package notification

import "time"

// ユーザーのFCM通知設定を管理する構造体
type NotificationSetting struct {
	ID        uint       `json:"id" gorm:"primarykey"`
	CreatedAt time.Time  `json:"createdAt"`                        // 作成時に自動設定される
	UpdatedAt time.Time  `json:"updatedAt"`                        // 更新時に自動設定される
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"index"` // ソフトデリート用
	UserID    string     `json:"userId" gorm:"not null;uniqueIndex:idx_user_id"`
	IsEnabled bool       `json:"isEnabled" gorm:"default:true"`
	FcmToken  string     `json:"fcmToken" gorm:"not null;uniqueIndex:idx_fcm_token"`
	Platform  string     `json:"platform" gorm:"not null;uniqueIndex:idx_platform"`
}
