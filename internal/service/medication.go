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
	
	// 日付ごとに整理したログを取得（同じ日の重複を除去）
	dateLogMap := make(map[string]model.MedicationLog)
	for _, log := range logs {
		dateStr := log.CreatedAt.Format("2006-01-02")
		// 同じ日付の場合は最新のログを使用
		// logsは新しい順にソートされているため、最初に見つけたログが最新
		if _, exists := dateLogMap[dateStr]; !exists {
			dateLogMap[dateStr] = log
		}
	}
	
	// 日付を過去順（降順）にソート
	var dates []time.Time
	for _, log := range dateLogMap {
		dates = append(dates, log.CreatedAt)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].After(dates[j])
	})
	
	// 連続出血日数を計算
	consecutiveBleedingDays := 0
	var consecutiveBleedingDates []time.Time
	lastDate := time.Time{}
	
	for _, date := range dates {
		// 日付を正規化（時間部分を削除）
		currDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		
		// 日付の文字列を取得
		dateStr := currDate.Format("2006-01-02")
		log := dateLogMap[dateStr]
		
		if log.HasBleeding {
			// 初めての出血日または連続している場合
			if consecutiveBleedingDays == 0 || lastDate.IsZero() {
				consecutiveBleedingDays = 1
				consecutiveBleedingDates = append(consecutiveBleedingDates, currDate)
			} else {
				// 日付の差を計算
				dayDiff := int(lastDate.Sub(currDate).Hours() / 24)
				
				// 前日からの連続か確認
				if dayDiff == 1 {
					consecutiveBleedingDays++
					consecutiveBleedingDates = append(consecutiveBleedingDates, currDate)
				} else {
					// 日付が連続していない場合はリセット
					consecutiveBleedingDays = 1
					consecutiveBleedingDates = []time.Time{currDate}
				}
			}
		} else {
			// 出血がない場合はリセット
			break
		}
		
		lastDate = currDate
	}
	
	// 連続3日間以上の出血がある場合、休薬期間判定
	if consecutiveBleedingDays >= 3 && len(consecutiveBleedingDates) >= 3 {
		// 休薬開始日は連続出血の最初の日
		restStartDate := consecutiveBleedingDates[len(consecutiveBleedingDates)-1]
		
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
	
	// 日付ごとに整理したログを取得（同じ日の重複を除去）
	dateMap := make(map[string]time.Time)
	for _, log := range logs {
		dateStr := log.CreatedAt.Format("2006-01-02")
		// 既に同じ日付のログがある場合は最新のものを使用
		if _, exists := dateMap[dateStr]; !exists {
			dateMap[dateStr] = log.CreatedAt
		}
	}
	
	// 日付を過去順（降順）にソート
	var dates []time.Time
	for _, date := range dateMap {
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].After(dates[j])
	})
	
	// 現在の日付の前日
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)
	
	// 休薬期間後からの連続日数をカウント
	streak := 0
	lastDate := today // 最初は今日から開始
	
	for _, date := range dates {
		// 日付を正規化（時間部分を削除）
		currDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		
		// 休薬期間後のログのみカウント
		if !lastRestPeriodEndDate.IsZero() && currDate.Before(lastRestPeriodEndDate) {
			break
		}
		
		// 現在の日付または前日の場合はカウント
		if currDate.Equal(today) || currDate.Equal(yesterday) {
			streak++
			lastDate = currDate
			continue
		}
		
		// 日付の差を計算
		dayDiff := int(lastDate.Sub(currDate).Hours() / 24)
		
		// 日付が連続しているか確認
		if dayDiff == 1 {
			// 前日のログであればカウント
			streak++
			lastDate = currDate
		} else {
			// 日付が連続していない場合は終了
			break
		}
	}
	
	return streak
}
