package user

type UserService struct {
	// UserRepository ユーザーデータへのアクセスを提供するリポジトリ
	UserRepository *UserRepository
}

// NewUserService は新しいUserServiceのインスタンスを作成する
func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

// GetAllUsers は全ユーザー情報を取得する
func (us *UserService) GetAllUsers() ([]User, error) {
	users, err := us.UserRepository.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}
