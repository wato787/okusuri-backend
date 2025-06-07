package service

import (
	"okusuri-backend/internal/dto"
	"okusuri-backend/internal/model"
	"okusuri-backend/internal/repository"
	"sort"
	"time"
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
		CurrentStreak:           0,
		IsRestPeriod:            false,
		RestDaysLeft:            0,
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
		restEndDate = time.Date(
			restEndDate.Year(), restEndDate.Month(), restEndDate.Day(),
			23, 59, 59, 0, restEndDate.Location(),
		)

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
	lastRestPeriodEndDate := s.findLastRestPeriodEndDate(logs)
	uniqueDates := s.extractUniqueDates(logs)
	return s.countConsecutiveDays(uniqueDates, lastRestPeriodEndDate, now)
}

// findLastRestPeriodEndDate は最後の休薬期間終了日を探す
func (s *MedicationService) findLastRestPeriodEndDate(logs []model.MedicationLog) time.Time {
	consecutiveBleedingCount := 0
	var bleedingDates []time.Time

	for i, log := range logs {
		if log.HasBleeding {
			consecutiveBleedingCount++
			bleedingDates = append(bleedingDates, log.CreatedAt)

			if consecutiveBleedingCount >= 3 {
				oldestBleedingDate := bleedingDates[len(bleedingDates)-1]
				return oldestBleedingDate.AddDate(0, 0, 4)
			}
		} else {
			consecutiveBleedingCount = 0
			bleedingDates = nil

			if i >= 3 && logs[i-1].HasBleeding && logs[i-2].HasBleeding && logs[i-3].HasBleeding {
				return log.CreatedAt
			}
		}
	}
	return time.Time{}
}

// extractUniqueDates は重複を除去した日付リストを取得する
func (s *MedicationService) extractUniqueDates(logs []model.MedicationLog) []time.Time {
	dateMap := make(map[string]time.Time)
	for _, log := range logs {
		dateStr := log.CreatedAt.Format("2006-01-02")
		if _, exists := dateMap[dateStr]; !exists {
			dateMap[dateStr] = log.CreatedAt
		}
	}

	var dates []time.Time
	for _, date := range dateMap {
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].After(dates[j])
	})
	return dates
}

// countConsecutiveDays は休薬期間後からの連続日数をカウントする
func (s *MedicationService) countConsecutiveDays(dates []time.Time, restEndDate time.Time, now time.Time) int {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)

	streak := 0
	lastDate := today

	for _, date := range dates {
		currDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

		if !restEndDate.IsZero() && currDate.Before(restEndDate) {
			break
		}

		if currDate.Equal(today) || currDate.Equal(yesterday) {
			streak++
			lastDate = currDate
			continue
		}

		dayDiff := int(lastDate.Sub(currDate).Hours() / 24)
		if dayDiff == 1 {
			streak++
			lastDate = currDate
		} else {
			break
		}
	}
	return streak
}

// GetMedicationStats はユーザーの服薬統計を計算する
func (s *MedicationService) GetMedicationStats(userID string) (*dto.MedicationStatsResponse, error) {
	// 統計データを取得
	logs, err := s.medicationRepo.GetStatsData(userID)
	if err != nil {
		return nil, err
	}

	// 最長連続服用日数を計算
	longestStreak := s.calculateLongestContinuousDays(logs)

	// 総出血日数を計算
	totalBleedingDays := s.calculateTotalBleedingDays(logs)

	// 服用中断回数を計算
	medicationBreaks := s.calculateMedicationBreaks(logs)

	// 平均サイクル長を計算
	averageCycleLength := s.calculateAverageCycleLength(logs)

	// 月間服用率を計算
	monthlyRate := s.calculateMonthlyMedicationRate(logs)

	return &dto.MedicationStatsResponse{
		LongestContinuousDays:  longestStreak,
		TotalBleedingDays:     totalBleedingDays,
		MedicationBreaks:      medicationBreaks,
		AverageCycleLength:    averageCycleLength,
		MonthlyMedicationRate: monthlyRate,
	}, nil
}

// calculateLongestContinuousDays は最長連続服用日数を計算する
func (s *MedicationService) calculateLongestContinuousDays(logs []model.MedicationLog) int {
	if len(logs) == 0 {
		return 0
	}

	// 日付ごとに整理（重複除去）
	dateMap := make(map[string]bool)
	for _, log := range logs {
		dateStr := log.CreatedAt.Format("2006-01-02")
		dateMap[dateStr] = true
	}

	// 日付をソート
	var dates []time.Time
	for dateStr := range dateMap {
		date, _ := time.Parse("2006-01-02", dateStr)
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	// 最長連続日数を計算
	maxStreak := 1
	currentStreak := 1

	for i := 1; i < len(dates); i++ {
		dayDiff := int(dates[i].Sub(dates[i-1]).Hours() / 24)
		if dayDiff == 1 {
			currentStreak++
			if currentStreak > maxStreak {
				maxStreak = currentStreak
			}
		} else {
			currentStreak = 1
		}
	}

	return maxStreak
}

// calculateTotalBleedingDays は総出血日数を計算する
func (s *MedicationService) calculateTotalBleedingDays(logs []model.MedicationLog) int {
	bleedingDates := make(map[string]bool)
	for _, log := range logs {
		if log.HasBleeding {
			dateStr := log.CreatedAt.Format("2006-01-02")
			bleedingDates[dateStr] = true
		}
	}
	return len(bleedingDates)
}

// calculateMedicationBreaks は服用中断回数を計算する
func (s *MedicationService) calculateMedicationBreaks(logs []model.MedicationLog) int {
	if len(logs) <= 1 {
		return 0
	}

	// 日付ごとに整理（重複除去）
	dateMap := make(map[string]bool)
	for _, log := range logs {
		dateStr := log.CreatedAt.Format("2006-01-02")
		dateMap[dateStr] = true
	}

	// 日付をソート
	var dates []time.Time
	for dateStr := range dateMap {
		date, _ := time.Parse("2006-01-02", dateStr)
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	// 中断回数をカウント
	breaks := 0
	for i := 1; i < len(dates); i++ {
		dayDiff := int(dates[i].Sub(dates[i-1]).Hours() / 24)
		if dayDiff > 1 {
			breaks++
		}
	}

	return breaks
}

// calculateAverageCycleLength は平均サイクル長を計算する
func (s *MedicationService) calculateAverageCycleLength(logs []model.MedicationLog) float64 {
	if len(logs) < 2 {
		return 0.0
	}

	// 服用期間を計算
	firstDate := logs[0].CreatedAt
	lastDate := logs[len(logs)-1].CreatedAt
	totalDays := int(lastDate.Sub(firstDate).Hours()/24) + 1

	// 日付ごとに整理（重複除去）
	dateMap := make(map[string]bool)
	for _, log := range logs {
		dateStr := log.CreatedAt.Format("2006-01-02")
		dateMap[dateStr] = true
	}

	medicationDays := len(dateMap)
	if medicationDays == 0 {
		return 0.0
	}

	// 平均サイクル長 = 総日数 / 服用日数
	return float64(totalDays) / float64(medicationDays) * 21.0 // 21日サイクルを基準
}

// calculateMonthlyMedicationRate は月間服用率を計算する
func (s *MedicationService) calculateMonthlyMedicationRate(logs []model.MedicationLog) float64 {
	if len(logs) == 0 {
		return 0.0
	}

	// 直近30日間の服用率を計算
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	medicationDays := 0
	dateMap := make(map[string]bool)

	for _, log := range logs {
		if log.CreatedAt.After(thirtyDaysAgo) {
			dateStr := log.CreatedAt.Format("2006-01-02")
			if !dateMap[dateStr] {
				dateMap[dateStr] = true
				medicationDays++
			}
		}
	}

	return float64(medicationDays) / 30.0 * 100.0
}
