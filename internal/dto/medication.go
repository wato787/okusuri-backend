package dto

import "time"

// 服用記録リクエスト
type MedicationLogRequest struct {
	HasBleeding bool       `json:"hasBleeding"`
	Date        *time.Time `json:"date,omitempty"` // 指定された日付（省略時は現在日時）
}
