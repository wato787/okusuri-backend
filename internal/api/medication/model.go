package medication

import "time"

// 服用履歴の構造体
type MedicationLog struct {
	ID          uint       `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time  `json:"createdAt"`                        // 作成時に自動設定される
	UpdatedAt   time.Time  `json:"updatedAt"`                        // 更新時に自動設定される
	DeletedAt   *time.Time `json:"deletedAt,omitempty" gorm:"index"` // ソフトデリート用
	UserID      string     `json:"userId" gorm:"not null;uniqueIndex:idx_user_id"`
	HasBleeding bool       `json:"hasBleeding" gorm:"default:false"`
}
