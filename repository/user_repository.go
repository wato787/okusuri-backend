package repository

// UserRepository はユーザーデータへのアクセスを提供する
type UserRepository struct{}

// NewUserRepository は新しいUserRepositoryのインスタンスを作成する
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}
