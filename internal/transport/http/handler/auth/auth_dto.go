package auth

type SignUpRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type SignInRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ErrorResponse представляет структуру ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error" example:"описание ошибки"`
}

// Response представляет структуру успешного ответа
type Response struct {
	Message string `json:"message" example:"операция выполнена успешно"`
}

// LogoutRequest представляет структуру запроса для выхода
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIs..."`
}
