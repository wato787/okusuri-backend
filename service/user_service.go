package service

import (
	"fmt"
	"okusuri-backend/model"
	"okusuri-backend/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// ユーザーをプロバイダーIDで取得する
func (s *UserService) GetUserByProviderId(providerId string) (*model.User, error) {
	user, err := s.userRepo.FindByProviderId(providerId)
	if err != nil {
		return nil, fmt.Errorf("ユーザーの取得に失敗しました: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("ユーザーが見つかりません")
	}

	return user, nil
}
