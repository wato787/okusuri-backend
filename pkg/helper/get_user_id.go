package helper

import (
	"errors"
	"okusuri-backend/internal/model"

	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(c *gin.Context) (string, error) {
	userInterface, exists := c.Get("user")
	if !exists {
		return "", errors.New("ユーザー情報がコンテキストに存在しません")
	}

	user, ok := userInterface.(*model.User)
	if !ok {
		return "", errors.New("ユーザー情報の型変換に失敗しました")
	}

	return user.ID, nil
}
