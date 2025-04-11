package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// DB接続
	setupDB()

	// Ginのルーターを作成
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // 0.0.0.0:8080 でサーバーを立てます。
}


func setupDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("環境変数の読み込みに失敗しました")
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URLが設定されていません")
	}

	db,err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("DB接続に失敗しました")
	}

	log.Println("DB接続に成功しました",db.Name())

}