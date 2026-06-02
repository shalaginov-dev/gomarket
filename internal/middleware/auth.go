package middleware

import (
	"gomarket/internal/token"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// 1. Проверяем наличие заголовка
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		// 2. Проверяем формат "Bearer <token>"
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// 3. Извлекаем токен
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 4. Валидируем токен (вызываем экспортированную функцию)
		claims, err := token.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		// 5. Сохраняем user_id в контексте Gin
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
