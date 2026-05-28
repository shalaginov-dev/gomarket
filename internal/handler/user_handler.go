package handler

import (
	"gomarket/internal/cache"
	"gomarket/internal/service"
	"gomarket/internal/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService   *service.UserService
	tokenStore    *cache.TokenStore
	jwtSecret     string
	jwtExpiry     int
	refreshExpiry int
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

func NewUserHandler(userService *service.UserService, tokenStore *cache.TokenStore, JWTSecret string, JWTExpiry int, refreshExpiry int) *UserHandler {
	return &UserHandler{
		userService:   userService,
		tokenStore:    tokenStore,
		jwtSecret:     JWTSecret,
		jwtExpiry:     JWTExpiry,
		refreshExpiry: refreshExpiry,
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
	accessToken, err := token.GenerateAccessToken(user.UserID, user.Email, user.Role, r.jwtSecret, r.jwtExpiry)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate access token"})
		return
	}
	refreshToken, err := token.GenerateRefreshToken(user.UserID, user.Role, r.jwtSecret, r.refreshExpiry)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate refresh token"})
		return
	}
	if err := r.tokenStore.SaveRefreshToken(c, user.UserID, refreshToken, time.Duration(r.refreshExpiry)*time.Minute); err != nil {
		c.JSON(500, gin.H{"error": "failed to save refresh token"})
		return
	}

	resUser := UserResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
	c.JSON(200, gin.H{
		"message":       "successfully logined",
		"user":          resUser,
		"token":         accessToken,
		"refresh_token": refreshToken,
	})
}
func (r *UserHandler) LogoutHandler(c *gin.Context) {
	value, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	userID, ok := value.(int)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid user_id"})
		return
	}
	if err := r.tokenStore.DeleteRefreshToken(c, userID); err != nil {
		c.JSON(500, gin.H{"error": "failed to logout"})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully logout",
	})
}
