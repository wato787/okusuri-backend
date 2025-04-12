package service

import (
	"fmt"
	"okusuri-backend/dto"
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

// ユーザーを登録する
func (s *UserService) RegisterUser(req dto.SignupRequest) (*model.User, error) {

	// ユーザーの存在確認
	existingUser, err := s.userRepo.FindByProviderId(req.ProviderID)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, nil // ユーザーは既に存在する
	}

	// 新しいユーザーの作成
	newUser := &model.User{
		Email:      req.Email,
		Name:       req.Name,
		ImageUrl:   req.ImageURL,
		Provider:   model.ProviderGoogle,
		ProviderId: req.ProviderID,
	}

	// ユーザーの保存
	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil

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
