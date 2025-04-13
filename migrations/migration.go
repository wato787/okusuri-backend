package migrations

import (
	"log"

	"gorm.io/gorm"
)

// RunMigrations はデータベースマイグレーションを実行します
func RunMigrations(db *gorm.DB) {
	log.Println("マイグレーションを実行します...")

	// // マイグレーション対象のモデルをここに追加
	err := db.AutoMigrate()
	if err != nil {
		log.Fatalf("マイグレーションに失敗しました: %v", err)
	}

	log.Println("マイグレーションが正常に完了しました")
}
