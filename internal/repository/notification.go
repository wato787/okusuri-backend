package repository

import (
	"okusuri-backend/internal/model"
	"okusuri-backend/pkg/config"
)

type NotificationRepository struct{}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

// GetSettingByUserID はユーザーIDに基づいて通知設定を取得する
func (r *NotificationRepository) GetSettingByUserID(userID string) (*model.NotificationSetting, error) {
	// DB接続
	db := config.DB

	// ユーザーIDに基づいて通知設定を取得
	var setting model.NotificationSetting
	if err := db.Where("user_id = ?", userID).First(&setting).Error; err != nil {
		return nil, err
	}

	return &setting, nil
}

// RegisterSetting は通知設定を登録する
func (r *NotificationRepository) RegisterSetting(setting *model.NotificationSetting) error {
	// DB接続
	db := config.DB

	if err := db.Create(setting).Error; err != nil {
		return err
	}

	return nil
}

// GetAllSettings は全ての通知設定を取得する
func (r *NotificationRepository) GetAllSettings() ([]model.NotificationSetting, error) {
	// DB接続
	db := config.DB

	var settings []model.NotificationSetting
	if err := db.Find(&settings).Error; err != nil {
		return nil, err
	}

	return settings, nil
}
