package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger はリクエストのログを記録するミドルウェアです
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// リクエスト処理の前
		c.Next()

		// リクエスト処理の後
		latency := time.Since(start)
		status := c.Writer.Status()
		log.Printf("| %3d | %13v | %s", status, latency, path)
	}
}

// CORS ヘッダーを設定するミドルウェア
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// TODO:tokenからユーザー情報を取得するミドルウェア作成
