package user

import (
	"net/http"
	"strconv"

	userDB "github.com/NikitaBelov-mobile/car-social/internal/database/user"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userRepo userDB.UserRepository
}

func NewHandler(userRepo userDB.UserRepository) *Handler {
	return &Handler{
		userRepo: userRepo,
	}
}

func (h *Handler) Register(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("/:id", h.getByID) // Получение пользователя по ID
		users.PUT("/:id", h.update)  // Обновление данных пользователя
	}
}

func (h *Handler) create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, существует ли пользователь
	existingUser, err := h.userRepo.GetByPhone(req.Phone)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	// Создаем нового пользователя
	user := &userDB.User{
		Phone:        req.Phone,
		PasswordHash: string(hashedPassword),
	}

	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, Response{
		ID:        user.ID,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

// GetByID godoc
// @Summary Получение пользователя
// @Tags users
// @Description Получение информации о пользователе по ID
// @Accept  json
// @Produce  json
// @Param id path int true "ID пользователя"
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse "неверный формат ID"
// @Failure 404 {object} ErrorResponse "пользователь не найден"
// @Failure 500 {object} ErrorResponse "внутренняя ошибка сервера"
// @Router /users/{id} [get]
func (h *Handler) getByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	user, err := h.userRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, Response{
		ID:        user.ID,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

// Update godoc
// @Summary Обновление пользователя
// @Tags users
// @Description Обновление данных пользователя
// @Accept  json
// @Produce  json
// @Param id path int true "ID пользователя"
// @Param input body UpdateRequest true "Данные для обновления"
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse "неверный формат данных"
// @Failure 404 {object} ErrorResponse "пользователь не найден"
// @Failure 409 {object} ErrorResponse "телефон уже занят"
// @Failure 500 {object} ErrorResponse "внутренняя ошибка сервера"
// @Router /users/{id} [put]
func (h *Handler) update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем существующего пользователя
	user, err := h.userRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Обновляем только переданные поля
	if req.Phone != "" {
		// Проверяем, не занят ли телефон другим пользователем
		existingUser, err := h.userRepo.GetByPhone(req.Phone)
		if err == nil && existingUser != nil && existingUser.ID != id {
			c.JSON(http.StatusConflict, gin.H{"error": "phone number already taken"})
			return
		}
		user.Phone = req.Phone
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
			return
		}
		user.PasswordHash = string(hashedPassword)
	}

	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, Response{
		ID:        user.ID,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}
