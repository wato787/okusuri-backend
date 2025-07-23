package repository

import (
	"fmt"
	"okusuri-backend/internal/model"
	"okusuri-backend/pkg/config"
	"time"
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

// GetConsecutiveDays はユーザーの連続服薬日数を計算する
func (r *MedicationRepository) GetConsecutiveDays(userID string) (int, error) {
	db := config.DB

	// ユーザーの服薬ログを日付降順で取得
	var logs []model.MedicationLog
	if err := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return 0, err
	}

	if len(logs) == 0 {
		return 0, nil
	}

	// 連続日数をカウント
	consecutiveDays := 1
	today := time.Now().Truncate(24 * time.Hour)

	// 最新の記録が今日かどうかチェック
	latestLog := logs[0]
	latestDate := latestLog.CreatedAt.Truncate(24 * time.Hour)

	// 最新の記録が今日でない場合は連続日数は0
	if !latestDate.Equal(today) {
		return 0, nil
	}

	// 前日から遡って連続日数をカウント
	expectedDate := today.AddDate(0, 0, -1)

	for i := 1; i < len(logs); i++ {
		logDate := logs[i].CreatedAt.Truncate(24 * time.Hour)

		if logDate.Equal(expectedDate) {
			consecutiveDays++
			expectedDate = expectedDate.AddDate(0, 0, -1)
		} else {
			// 連続していない場合は終了
			break
		}
	}

	return consecutiveDays, nil
}
