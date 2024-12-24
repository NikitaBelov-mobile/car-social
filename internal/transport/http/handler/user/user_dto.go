package user

type CreateRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateRequest представляет структуру запроса на обновление пользователя
type UpdateRequest struct {
	Phone    string `json:"phone,omitempty" example:"79991234567"`
	Password string `json:"password,omitempty" example:"newpassword123" binding:"omitempty,min=6"`
}

// Response представляет структуру ответа с данными пользователя
type Response struct {
	ID        int    `json:"id" example:"1"`
	Phone     string `json:"phone" example:"79991234567"`
	CreatedAt string `json:"created_at" example:"2024-03-20 15:04:05"`
}

// ErrorResponse представляет структуру ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error" example:"описание ошибки"`
}
