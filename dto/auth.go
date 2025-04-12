package dto

// 新規登録のリクエストDTO
type SignupRequest struct {
	IDToken    string `json:"id_token" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Name       string `json:"name" binding:"required"`
	ImageURL   string `json:"image_url"` // 任意フィールド
	ProviderID string `json:"provider_id" binding:"required"`
}

// AuthResponse は認証レスポンスのDTOを定義します
type AuthResponse struct {
	User      UserResponse `json:"user"`
	Token     string       `json:"token"`
	ExpiresAt int64        `json:"expires_at,omitempty"`
	Message   string       `json:"message,omitempty"`
}
