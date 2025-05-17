package service

import (
	"sort"
	"time"
	"okusuri-backend/internal/dto"
	"okusuri-backend/internal/model"
	"okusuri-backend/internal/repository"
)

type MedicationService struct {
	medicationRepo *repository.MedicationRepository
}

func NewMedicationService(medicationRepo *repository.MedicationRepository) *MedicationService {
	return &MedicationService{
		medicationRepo: medicationRepo,
	}
}

// GetMedicationStatus は現在の服薬ステータスを計算する
func (s *MedicationService) GetMedicationStatus(userID string) (*dto.MedicationStatusResponse, error) {
	// 服薬ログを取得
	logs, err := s.medicationRepo.GetLogsByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 日付でソート（新しい順）
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].CreatedAt.After(logs[j].CreatedAt)
	})

	// 現在日時
	now := time.Now()

	// デフォルトのレスポンス
	response := &dto.MedicationStatusResponse{
		CurrentStreak:          0,
		IsRestPeriod:           false,
		RestDaysLeft:           0,
		ConsecutiveBleedingDays: 0,
	}

	// ログが存在しない場合は初期値を返す
	if len(logs) == 0 {
		return response, nil
	}

	// 休薬期間の判定と連続出血日数の計算
	isInRestPeriod, restDaysLeft, consecutiveBleedingDays := s.calculateRestPeriodStatus(logs, now)
	response.IsRestPeriod = isInRestPeriod
	response.RestDaysLeft = restDaysLeft
	response.ConsecutiveBleedingDays = consecutiveBleedingDays

	// 休薬期間中でなければ、現在の連続服用日数を計算
	if !isInRestPeriod {
		response.CurrentStreak = s.calculateCurrentStreak(logs, now)
	}

	return response, nil
}

// calculateRestPeriodStatus は休薬期間の状態を計算する
func (s *MedicationService) calculateRestPeriodStatus(logs []model.MedicationLog, now time.Time) (bool, int, int) {
	const restPeriodDays = 4 // 休薬期間は4日間
	
	// 最初に連続出血日数を計算
	consecutiveBleedingDays := 0
	for _, log := range logs {
		if log.HasBleeding {
			consecutiveBleedingDays++
		} else {
			break
		}
	}

	// 連続3日間以上の出血がある場合、休薬期間判定
	if consecutiveBleedingDays >= 3 {
		// 最初の出血記録の日付から休薬期間を計算
		// 記録は新しい順に並んでいるため、連続3日目の記録はインデックス2
		if len(logs) >= 3 {
			// 休薬開始日は3日連続出血の最初の日
			restStartDate := logs[consecutiveBleedingDays-1].CreatedAt
			
			// 休薬終了日は休薬開始日から4日後の終日
			restEndDate := restStartDate.AddDate(0, 0, restPeriodDays)
			restEndDate = time.Date(restEndDate.Year(), restEndDate.Month(), restEndDate.Day(), 23, 59, 59, 0, restEndDate.Location())
			
			// 現在が休薬期間内かどうか
			if now.Before(restEndDate) {
				// 残り日数を計算（日単位で切り上げ）
				duration := restEndDate.Sub(now)
				daysLeft := int(duration.Hours() / 24)
				if duration.Hours() > float64(daysLeft*24) {
					daysLeft++
				}
				return true, daysLeft, consecutiveBleedingDays
			}
		}
	}

	return false, 0, consecutiveBleedingDays
}

// calculateCurrentStreak は現在の連続服用日数を計算する
func (s *MedicationService) calculateCurrentStreak(logs []model.MedicationLog, now time.Time) int {
	// 最後の休薬期間を探す
	lastRestPeriodEndDate := time.Time{}
	
	// 3日連続出血を探す（新しい順からチェック）
	consecutiveBleedingCount := 0
	var bleedingDates []time.Time
	
	for i, log := range logs {
		if log.HasBleeding {
			consecutiveBleedingCount++
			bleedingDates = append(bleedingDates, log.CreatedAt)
			
			// 3日連続出血を見つけた場合
			if consecutiveBleedingCount >= 3 {
				// 最も古い連続出血日から4日後を休薬期間終了日とする
				oldestBleedingDate := bleedingDates[len(bleedingDates)-1]
				lastRestPeriodEndDate = oldestBleedingDate.AddDate(0, 0, 4)
				break
			}
		} else {
			// 連続が途切れたらリセット
			consecutiveBleedingCount = 0
			bleedingDates = nil
			
			// 非出血ログが見つかった場合、前の3日間をチェック
			if i >= 3 && logs[i-1].HasBleeding && logs[i-2].HasBleeding && logs[i-3].HasBleeding {
				// 3日連続出血の後の初めての非出血ログを見つけた
				lastRestPeriodEndDate = log.CreatedAt
				break
			}
		}
	}
	
	// 休薬期間後の最初のログから数える
	streak := 0
	if !lastRestPeriodEndDate.IsZero() {
		// 休薬期間後のログをカウント
		for _, log := range logs {
			if log.CreatedAt.After(lastRestPeriodEndDate) {
				streak++
			} else {
				break
			}
		}
	} else {
		// 休薬期間がない場合は、すべてのログをカウント（日付ベースでユニーク）
		dateSet := make(map[string]bool)
		for _, log := range logs {
			dateStr := log.CreatedAt.Format("2006-01-02")
			dateSet[dateStr] = true
		}
		streak = len(dateSet)
	}
	
	return streak
}
