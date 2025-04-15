package model

import "gorm.io/gorm"

// ユーザーのFCMトークンを管理する構造体
type FcmToken struct {
	gorm.Model
	Token  string `gorm:"unique" json:"token"` // FCMトークン
	UserId uint   `json:"user_id"`             // ユーザーID
}
