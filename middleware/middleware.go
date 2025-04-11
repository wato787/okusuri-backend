package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
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

// JWTAuth は認証を処理するミドルウェア
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		// TODO トークンの検証ロジックをここに実装
		if token == "" {
			c.JSON(401, gin.H{"error": "認証が必要です"})
			c.Abort()
			return
		}

		// 検証が成功した場合、ユーザー情報をコンテキストに設定
		// c.Set("userId", userId)

		c.Next()
	}
}
