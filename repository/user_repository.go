package repository

import (
	"okusuri-backend/config"
	"okusuri-backend/model"
)

// UserRepository はユーザーデータへのアクセスを提供する
type UserRepository struct{}

// NewUserRepository は新しいUserRepositoryのインスタンスを作成する
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// FindByProviderId は指定されたプロバイダーIDを持つユーザーを検索する
func (repo *UserRepository) FindByProviderId(providerId string) (*model.User, error) {
	var user model.User
	result := config.DB.Where("provider_id = ?", providerId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil // ユーザーが見つからない場合
	}
	return &user, nil
}

// Create は新しいユーザーをデータベースに保存する
func (repo *UserRepository) Create(user *model.User) error {
	result := config.DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
