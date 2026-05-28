package main

import (
	"context"
	"fmt"
	"log"

	"gomarket/internal/cache"
	"gomarket/internal/config"
	"gomarket/internal/db"
	"gomarket/internal/handler"
	"gomarket/internal/middleware"
	"gomarket/internal/repository"
	"gomarket/internal/service"

	"github.com/gin-gonic/gin"
)

type Health struct {
	Status string `json:"status"`
}

func healthHandler(c *gin.Context) {
	health := Health{Status: "ok"}

	c.JSON(200, health)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Создаем пул соединений с БД
	pool, err := db.NewPool(context.Background(), cfg.DBDSN)
	if err != nil {
		log.Fatal(err)
	}
	// Закрываем соединение с базой, когда приложение останавливается
	defer pool.Close()

	redis, err := cache.NewRedisClient(context.Background(), cfg.RedisURL)
	if err != nil {
		log.Fatal(err)
	}
	defer redis.Close()

	// Делаем миграцию
	if err := db.RunMigrations(cfg.DBDSN, "./migrations"); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(pool)
	passwordService := service.NewPasswordServise()
	tokenStore := cache.NewTokenStore(redis)
	userService := service.NewUserService(userRepo, passwordService)
	userHandler := handler.NewUserHandler(userService, tokenStore, cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshExpiry)

	// Создаем роутер
	r := gin.Default()
	// Добавляем эндпоинты
	r.GET("/health", middleware.AuthMiddleware(cfg.JWTSecret), healthHandler)
	r.POST("/register", userHandler.RegisterHandler)
	r.POST("/login", userHandler.LoginHandler)
	r.POST("/logout", middleware.AuthMiddleware(cfg.JWTSecret), userHandler.LogoutHandler)
	r.POST("/refresh", userHandler.RefreshHandler)

	// Стартуем сервер
	if err := r.Run(":" + cfg.Port); err != nil {
		fmt.Println("Server failed:", err)
	}
}
