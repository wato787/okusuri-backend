package user

type Service struct {
	Repository *Repository
}

// NewUserService は新しいUserServiceのインスタンスを作成する
func NewService(repository *Repository) *Service {
	return &Service{
		Repository: repository,
	}
}

// GetAllUsers は全ユーザー情報を取得する
func (us *Service) GetAllUsers() ([]User, error) {
	users, err := us.Repository.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}
