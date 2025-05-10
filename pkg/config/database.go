package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB はアプリケーション全体で使用するデータベース接続です
var DB *gorm.DB

// SetupDB はデータベース接続を初期化します
func SetupDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not loaded")
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URLが設定されていません")
	}

	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("DB接続に失敗しました: ", err)
	}

	log.Println("DB接続に成功しました")
	DB = db
}

// GetDB はデータベース接続を返します
func GetDB() *gorm.DB {
	return DB
}
