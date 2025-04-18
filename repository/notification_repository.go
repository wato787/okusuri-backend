package repository

import (
	"okusuri-backend/config"
	"okusuri-backend/model"
)

type NotificationRepository struct{}

// NewNotificationRepository は新しいNotificationRepositoryのインスタンスを作成する
func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

// GetNotificationSettingByUserID はユーザーIDに基づいて通知設定を取得する
func (r *NotificationRepository) GetNotificationSettingByUserID(userID int) (*model.NotificationSetting, error) {
	// DB接続
	db := config.DB

	// ユーザーIDに基づいて通知設定を取得
	var notificationSetting model.NotificationSetting
	if err := db.Where("user_id = ?", userID).First(&notificationSetting).Error; err != nil {
		return nil, err
	}

	return &notificationSetting, nil
}
