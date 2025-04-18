package service

import (
	"okusuri-backend/model"
	"okusuri-backend/repository"
)

type UserService struct {
	// UserRepository ユーザーデータへのアクセスを提供するリポジトリ
	UserRepository *repository.UserRepository
}

// NewUserService は新しいUserServiceのインスタンスを作成する
func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

// GetAllUsers は全ユーザー情報を取得する
func (us *UserService) GetAllUsers() ([]model.User, error) {
	users, err := us.UserRepository.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}
