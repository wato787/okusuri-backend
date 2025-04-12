package dto

// エラーレスポンスの共通構造体
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
