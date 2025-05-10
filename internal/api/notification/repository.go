package notification

import (
	"okusuri-backend/pkg/config"
)

// Repository は通知設定のデータアクセスを提供する
type Repository struct{}

// NewRepository は新しいRepositoryインスタンスを作成する
func NewRepository() *Repository {
	return &Repository{}
}

// GetSettingByUserID はユーザーIDに基づいて通知設定を取得する
func (r *Repository) GetSettingByUserID(userID string) (*NotificationSetting, error) {
	// DB接続
	db := config.DB

	// ユーザーIDに基づいて通知設定を取得
	var setting NotificationSetting
	if err := db.Where("user_id = ?", userID).First(&setting).Error; err != nil {
		return nil, err
	}

	return &setting, nil
}

// RegisterSetting は通知設定を登録する
func (r *Repository) RegisterSetting(setting *NotificationSetting) error {
	// DB接続
	db := config.DB

	if err := db.Create(setting).Error; err != nil {
		return err
	}

	return nil
}

// GetAllSettings は全ての通知設定を取得する
func (r *Repository) GetAllSettings() ([]NotificationSetting, error) {
	// DB接続
	db := config.DB

	var settings []NotificationSetting
	if err := db.Find(&settings).Error; err != nil {
		return nil, err
	}

	return settings, nil
}
