package service

import (
	"okusuri-backend/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMedicationService_New(t *testing.T) {
	t.Run("MedicationServiceが正常に作成される", func(t *testing.T) {
		// MedicationServiceの作成をテスト
		assert.NotPanics(t, func() {
			NewMedicationService(nil)
		})
	})
}

// 統計計算ロジックのテスト
func TestMedicationService_CalculateLongestContinuousDays(t *testing.T) {
	service := &MedicationService{}

	t.Run("空のログの場合", func(t *testing.T) {
		logs := []model.MedicationLog{}
		result := service.calculateLongestContinuousDays(logs)
		assert.Equal(t, 0, result)
	})

	t.Run("連続した日付のログの場合", func(t *testing.T) {
		logs := []model.MedicationLog{
			{CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			{CreatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
			{CreatedAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
		}
		result := service.calculateLongestContinuousDays(logs)
		assert.Equal(t, 3, result)
	})

	t.Run("間隔のある日付のログの場合", func(t *testing.T) {
		logs := []model.MedicationLog{
			{CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			{CreatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
			{CreatedAt: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)}, // 間隔あり
			{CreatedAt: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)},
			{CreatedAt: time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC)},
		}
		result := service.calculateLongestContinuousDays(logs)
		assert.Equal(t, 3, result) // 最長は3日間
	})
}

func TestMedicationService_CalculateTotalBleedingDays(t *testing.T) {
	service := &MedicationService{}

	t.Run("出血日がない場合", func(t *testing.T) {
		logs := []model.MedicationLog{
			{CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), HasBleeding: false},
			{CreatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), HasBleeding: false},
		}
		result := service.calculateTotalBleedingDays(logs)
		assert.Equal(t, 0, result)
	})

	t.Run("出血日がある場合", func(t *testing.T) {
		logs := []model.MedicationLog{
			{CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), HasBleeding: true},
			{CreatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), HasBleeding: false},
			{CreatedAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), HasBleeding: true},
		}
		result := service.calculateTotalBleedingDays(logs)
		assert.Equal(t, 2, result)
	})

	t.Run("同じ日に複数の記録がある場合", func(t *testing.T) {
		logs := []model.MedicationLog{
			{CreatedAt: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), HasBleeding: true},
			{CreatedAt: time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC), HasBleeding: true}, // 同じ日
			{CreatedAt: time.Date(2024, 1, 2, 9, 0, 0, 0, time.UTC), HasBleeding: true},
		}
		result := service.calculateTotalBleedingDays(logs)
		assert.Equal(t, 2, result) // 重複は除去される
	})
}

func TestMedicationService_CalculateMedicationBreaks(t *testing.T) {
	service := &MedicationService{}

	t.Run("中断がない連続した記録の場合", func(t *testing.T) {
		logs := []model.MedicationLog{
			{CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			{CreatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
			{CreatedAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
		}
		result := service.calculateMedicationBreaks(logs)
		assert.Equal(t, 0, result)
	})

	t.Run("中断がある記録の場合", func(t *testing.T) {
		logs := []model.MedicationLog{
			{CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			{CreatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
			{CreatedAt: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)}, // 2日間の中断
			{CreatedAt: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)},
			{CreatedAt: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)}, // 3日間の中断
		}
		result := service.calculateMedicationBreaks(logs)
		assert.Equal(t, 2, result) // 2回の中断
	})
}
