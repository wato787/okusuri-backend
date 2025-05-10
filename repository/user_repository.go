package repository

import (
	"okusuri-backend/model"
	"okusuri-backend/pkg/config"
)

// UserRepository はユーザーデータへのアクセスを提供する
type UserRepository struct{}

// NewUserRepository は新しいUserRepositoryのインスタンスを作成する
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// tokenからユーザー情報を取得する
func (r *UserRepository) GetUserByToken(token string) (*model.User, error) {
	db := config.DB
	var user model.User
	var session model.Session

	// セッションを取得
	if err := db.Where("token = ?", token).First(&session).Error; err != nil {
		return nil, err
	}

	// ユーザーを取得
	if err := db.Where("id = ?", session.UserID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	db := config.DB
	var users []model.User

	// ユーザーを取得
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
