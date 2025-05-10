package medication

import (
	"okusuri-backend/pkg/config"
)

// Repository インターフェース定義
type Repository interface {
	RegisterLog(userID string, log MedicationLog) error
	GetLogsByUserID(userID string) ([]MedicationLog, error)
}

// MedicationRepository はRepository インターフェースの実装
type MedicationRepository struct{}

// NewRepository は新しいRepository実装インスタンスを作成する
func NewRepository() Repository {
	return &MedicationRepository{}
}

// RegisterLog はユーザーの服用記録をデータベースに登録する
func (r *MedicationRepository) RegisterLog(userID string, log MedicationLog) error {
	// DB接続
	db := config.DB

	// ユーザーIDに基づいて服用履歴を登録
	if err := db.Create(&log).Error; err != nil {
		return err
	}

	return nil
}

// GetLogsByUserID はユーザーIDに基づいて服用履歴をデータベースから取得する
func (r *MedicationRepository) GetLogsByUserID(userID string) ([]MedicationLog, error) {
	// DB接続
	db := config.DB

	// ユーザーIDに基づいて服用履歴を取得
	var logs []MedicationLog
	if err := db.Where("user_id = ?", userID).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}
