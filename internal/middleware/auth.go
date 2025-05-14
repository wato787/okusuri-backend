package middleware

import (
	"okusuri-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

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
