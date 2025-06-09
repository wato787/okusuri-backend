package dto

// 薬の統計レスポンス
type MedicationStatsResponse struct {
	// 最長連続服用日数
	LongestContinuousDays int `json:"longestContinuousDays"`
	// 総出血日数
	TotalBleedingDays int `json:"totalBleedingDays"`
	// 服用中断回数
	MedicationBreaks int `json:"medicationBreaks"`
	// 平均サイクル長（日数）
	AverageCycleLength float64 `json:"averageCycleLength"`
	// 月間服用率（%）
	MonthlyMedicationRate float64 `json:"monthlyMedicationRate"`
}
