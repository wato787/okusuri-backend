package main

import (
	"log"
	"okusuri-backend/config"
	"okusuri-backend/migrations"
	"okusuri-backend/routes"
)

func main() {

	// DB接続
	config.SetupDB()

	// マイグレーションの実行
	migrations.RunMigrations(config.GetDB())

	// Ginのルーターを作成
	router := routes.SetupRoutes()

	// サーバーを起動
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
