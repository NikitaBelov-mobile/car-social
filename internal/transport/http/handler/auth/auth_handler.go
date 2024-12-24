package auth

import (
	"fmt"
	"net/http"

	authDB "github.com/NikitaBelov-mobile/car-social/internal/database/auth"
	userDB "github.com/NikitaBelov-mobile/car-social/internal/database/user"
	"github.com/NikitaBelov-mobile/car-social/internal/service/token"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// type Handler struct {
// 	userRepo userDB.UserRepository
// 	authRepo AuthRepository
// }

type AuthRepository interface {
	CreateSession(session *authDB.Session) error
	GetSessionByRefreshToken(refreshToken string) (*authDB.Session, error)
	DeleteSession(refreshToken string) error
	DeleteUserSessions(userID int) error
}

type Handler struct {
	userRepo     userDB.UserRepository
	authRepo     AuthRepository
	tokenManager *token.TokenManager
}

func NewHandler(userRepo userDB.UserRepository, authRepo AuthRepository, tokenManager *token.TokenManager) *Handler {
	return &Handler{
		userRepo:     userRepo,
		authRepo:     authRepo,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Register(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refresh)
		auth.POST("/logout", h.logout)
	}
}

// SignUp godoc
// @Summary Регистрация пользователя
// @Tags auth
// @Description Создание нового пользователя в системе
// @Accept  json
// @Produce  json
// @Param input body SignUpRequest true "Данные для регистрации"
// @Success 201 {object} Response "успешная регистрация"
// @Failure 400 {object} ErrorResponse "неверный формат данных"
// @Failure 409 {object} ErrorResponse "пользователь уже существует"
// @Failure 500 {object} ErrorResponse "внутренняя ошибка сервера"
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := h.userRepo.GetByPhone(req.Phone)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	user := &userDB.User{
		Phone:        req.Phone,
		PasswordHash: string(hashedPassword),
	}

	if err := h.userRepo.Create(user); err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.Status(http.StatusCreated)
}

// SignIn godoc
// @Summary Вход в систему
// @Tags auth
// @Description Аутентификация пользователя
// @Accept  json
// @Produce  json
// @Param input body SignInRequest true "Данные для входа"
// @Success 200 {object} TokensResponse "токены доступа"
// @Failure 400 {object} ErrorResponse "неверный формат данных"
// @Failure 401 {object} ErrorResponse "неверные учетные данные"
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetByPhone(req.Phone)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Генерируем токены
	accessToken, err := h.tokenManager.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := h.tokenManager.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	session := &authDB.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
	}

	if err := h.authRepo.CreateSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	c.JSON(http.StatusOK, TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Refresh godoc
// @Summary Обновление токена
// @Tags auth
// @Description Обновление access token с помощью refresh token
// @Accept  json
// @Produce  json
// @Param input body RefreshRequest true "Refresh token"
// @Success 200 {object} TokensResponse "новые токены"
// @Failure 400 {object} ErrorResponse "неверный формат данных"
// @Failure 401 {object} ErrorResponse "невалидный refresh token"
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.authRepo.GetSessionByRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	if err := h.authRepo.DeleteSession(req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to refresh session"})
		return
	}

	// Генерируем новые токены
	accessToken, err := h.tokenManager.GenerateAccessToken(session.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := h.tokenManager.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	newSession := &authDB.Session{
		UserID:       session.UserID,
		RefreshToken: refreshToken,
	}

	if err := h.authRepo.CreateSession(newSession); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create new session"})
		return
	}

	c.JSON(http.StatusOK, TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Logout godoc
// @Summary Выход из системы
// @Tags auth
// @Description Завершение сессии пользователя
// @Accept  json
// @Produce  json
// @Param input body LogoutRequest true "Refresh token"
// @Success 200 {object} Response "успешный выход"
// @Failure 400 {object} ErrorResponse "неверный формат данных"
// @Failure 500 {object} ErrorResponse "ошибка сервера"
// @Router /auth/logout [post]
func (h *Handler) logout(c *gin.Context) {
	// Получаем refresh token из тела запроса
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token не предоставлен"})
		return
	}

	// Удаляем сессию
	if err := h.authRepo.DeleteSession(req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при выходе из системы"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "успешный выход из системы"})
}
