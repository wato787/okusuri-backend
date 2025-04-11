package model

import (
	"gorm.io/gorm"
)

// ユーザー情報を表す構造体。OAuth2前提のためパスワードは持たない。
type User struct {
	gorm.Model
	Email      string       `json:"email" gorm:"unique;not null"`
	Name       string       `json:"name" gorm:"not null"`
	ImageUrl   string       `json:"image_url"`
	Provider   AuthProvider `json:"provider" gorm:"not null;uniqueIndex:idx_provider_provider_id"`
	ProviderId string       `json:"provider_id" gorm:"not null;uniqueIndex:idx_provider_provider_id"`
}

type AuthProvider string

// プロバイダーの定数
const (
	ProviderGoogle AuthProvider = "google"
)
