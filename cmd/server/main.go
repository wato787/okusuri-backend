package main

import (
	"log"
	"os"

	routes "okusuri-backend/internal"
	"okusuri-backend/migrations"
	"okusuri-backend/pkg/config"
)

func main() {
	// DB接続
	config.SetupDB()

	// マイグレーションの実行
	migrations.RunMigrations(config.GetDB())

	// Ginのルーターを作成
	router := routes.SetupRoutes()

	// サーバーを起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
