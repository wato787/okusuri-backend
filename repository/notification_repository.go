package repository

import (
	"okusuri-backend/config"
	"okusuri-backend/model"
)

type NotificationRepository struct{}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

func (r *NotificationRepository) GetNotificationSettingByUserID(userID string) (*model.NotificationSetting, error) {
	// DB接続
	db := config.DB

	// ユーザーIDに基づいて通知設定を取得
	var notificationSetting model.NotificationSetting
	if err := db.Where("user_id = ?", userID).First(&notificationSetting).Error; err != nil {
		return nil, err
	}

	return &notificationSetting, nil
}
