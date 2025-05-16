package repository

import (
	"fmt"
	"okusuri-backend/internal/model"
	"okusuri-backend/pkg/config"
)

type MedicationRepository struct{}

func NewMedicationRepository() *MedicationRepository {
	return &MedicationRepository{}
}

// RegisterLog はユーザーの服用記録をデータベースに登録する
func (r *MedicationRepository) RegisterLog(userID string, log model.MedicationLog) error {
	// DB接続
	db := config.DB

	// ユーザーIDに基づいて服用履歴を登録
	if err := db.Create(&log).Error; err != nil {
		return err
	}

	return nil
}

// GetLogsByUserID はユーザーIDに基づいて服用履歴をデータベースから取得する
func (r *MedicationRepository) GetLogsByUserID(userID string) ([]model.MedicationLog, error) {
	// DB接続
	db := config.DB

	// ユーザーIDに基づいて服用履歴を取得
	var logs []model.MedicationLog
	if err := db.Where("user_id = ?", userID).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

// GetLogByID はIDに基づいて単一の服薬ログを取得する
func (r *MedicationRepository) GetLogByID(userID string, logID uint) (*model.MedicationLog, error) {
	// DB接続
	db := config.DB

	// ユーザーIDとログIDに基づいて服薬履歴を取得
	var log model.MedicationLog
	if err := db.Where("id = ? AND user_id = ?", logID, userID).First(&log).Error; err != nil {
		return nil, err
	}

	return &log, nil
}

// UpdateLog は指定されたIDの服薬ログを更新する
func (r *MedicationRepository) UpdateLog(userID string, logID uint, hasBleeding bool) error {
	// DB接続
	db := config.DB

	// ユーザーIDとログIDに基づいてログを更新
	result := db.Model(&model.MedicationLog{}).
		Where("id = ? AND user_id = ?", logID, userID).
		Update("has_bleeding", hasBleeding)

	if result.Error != nil {
		return result.Error
	}

	// 更新された行数が0の場合は、ログが見つからないエラーを返す
	if result.RowsAffected == 0 {
		return fmt.Errorf("log not found or user not authorized")
	}

	return nil
}
