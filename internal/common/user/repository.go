package user

import (
	"okusuri-backend/pkg/config"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

// tokenからユーザー情報を取得する
func (r *Repository) GetUserByToken(token string) (*User, error) {
	db := config.DB
	var user User
	var session Session

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

func (r *Repository) GetAllUsers() ([]User, error) {
	db := config.DB
	var users []User

	// ユーザーを取得
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
