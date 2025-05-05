package dto

// 服用記録リクエスト
type MedicationLogRequest struct {
	HasBleeding bool `json:"has_bleeding"`
}
