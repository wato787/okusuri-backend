package repository

import (
	"okusuri-backend/model"
	"okusuri-backend/pkg/config"
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

func (r *NotificationRepository) RegisterNotificationSetting(notificationSetting *model.NotificationSetting) error {
	// DB接続
	db := config.DB

	if err := db.Create(notificationSetting).Error; err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) GetAllNotificationSettings() ([]model.NotificationSetting, error) {
	// DB接続
	db := config.DB

	var notificationSettings []model.NotificationSetting
	if err := db.Find(&notificationSettings).Error; err != nil {
		return nil, err
	}

	return notificationSettings, nil
}
