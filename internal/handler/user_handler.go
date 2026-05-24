package handler

import (
	"gomarket/internal/service"
	"gomarket/internal/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
	jwtSecret   string
	jwtExpiry   int
}
type UserResponse struct {
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func NewUserHandler(userService *service.UserService, JWTSecret string, JWTExpiry int) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtSecret:   JWTSecret,
		jwtExpiry:   JWTExpiry,
	}
}

func (r *UserHandler) RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	// Проверяем тело запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Создаем учетную запись ползователя
	user, err := r.userService.Register(c, req.Email, req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resUser := UserResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
	c.JSON(200, gin.H{
		"message": "User registered successfully",
		"user":    resUser,
	})
}

func (r *UserHandler) LoginHandler(c *gin.Context) {
	var req LoginRequest
	// Проверяем тело запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := r.userService.Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	accessToken, err := token.GenerateAccessToken(user.UserID, user.Role, r.jwtSecret, r.jwtExpiry)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}
	resUser := UserResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
	c.JSON(200, gin.H{
		"message": "successfully logined",
		"user":    resUser,
		"token":   accessToken,
	})
}
