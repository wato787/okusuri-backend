package main

import (
	"okusuri-backend/config"
	"okusuri-backend/migrations"

	"github.com/gin-gonic/gin"
)

func main() {

	// DB接続
	config.SetupDB()

	// マイグレーションの実行
	migrations.RunMigrations(config.GetDB())

	// Ginのルーターを作成
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	r.Run() // 0.0.0.0:8080 でサーバーを立てます。
}
