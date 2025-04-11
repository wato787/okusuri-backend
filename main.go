package main

import (
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

	router.Run() // 0.0.0.0:8080 でサーバーを立てます。
}
