package migrations

import (
	"log"
	"okusuri-backend/internal/model"

	"gorm.io/gorm"
)

// RunMigrations はデータベースマイグレーションを実行します
func RunMigrations(db *gorm.DB) {
	log.Println("マイグレーションを実行します...")

	// マイグレーション対象のモデルをここに追加
	err := db.AutoMigrate(
		&model.User{},
		&model.Session{},
		&model.Account{},
		&model.Verification{},
		&model.NotificationSetting{}, 
		&model.MedicationLog{},
	)
	if err != nil {
		log.Fatalf("マイグレーションに失敗しました: %v", err)
	}

	log.Println("マイグレーションが正常に完了しました")
}
