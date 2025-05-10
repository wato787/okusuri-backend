package middleware

import (
	"log"
	"okusuri-backend/repository"
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

func Auth(userRepository *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bearerトークンを取得
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		token := authHeader[len("Bearer "):]
		if token == "" {
			c.JSON(401, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}
		// sessionテーブルのtokenと一致するレコードを取得
		user, err := userRepository.GetUserByToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// ユーザー情報をコンテキストに保存
		c.Set("user", user)

		c.Next()
	}
}
