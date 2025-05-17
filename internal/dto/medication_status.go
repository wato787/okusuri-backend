package dto

// 服薬ステータスレスポンス
type MedicationStatusResponse struct {
	CurrentStreak           int  `json:"currentStreak"`           // 現在の連続服用日数
	IsRestPeriod            bool `json:"isRestPeriod"`            // 休薬期間中かどうか
	RestDaysLeft            int  `json:"restDaysLeft"`            // 休薬期間の残り日数（休薬期間中の場合）
	ConsecutiveBleedingDays int  `json:"consecutiveBleedingDays"` // 連続出血日数
}
