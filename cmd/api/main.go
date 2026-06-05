package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gomarket/internal/cache"
	"gomarket/internal/config"
	"gomarket/internal/db"
	"gomarket/internal/handler"
	"gomarket/internal/logger"
	"gomarket/internal/middleware"
	"gomarket/internal/repository"
	"gomarket/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	// Создаем логер
	zapLogger, err := logger.New(cfg.Env)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем пул соединений с БД
	pool, err := db.NewPool(context.Background(), cfg.DBDSN)
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	// Закрываем соединение с базой, когда приложение останавливается
	defer pool.Close()

	redis, err := cache.NewRedisClient(context.Background(), cfg.RedisURL)
	if err != nil {
		zapLogger.Fatal("Failed to connect to redis", zap.Error(err))
	}
	defer redis.Close()

	// Делаем миграцию
	if err := db.RunMigrations(cfg.DBDSN, "./migrations"); err != nil {
		zapLogger.Fatal("Failed to run migrations", zap.Error(err))
	}

	userRepo := repository.NewUserRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	passwordService := service.NewPasswordServise()
	tokenStore := cache.NewTokenStore(redis)
	userService := service.NewUserService(userRepo, passwordService)
	userHandler := handler.NewUserHandler(userService, tokenStore, cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshExpiry)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)
	// Создаем роутер
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(zapLogger))
	// Добавляем эндпоинты
	r.GET("/health", middleware.AuthMiddleware(cfg.JWTSecret), healthHandler)
	r.POST("/register", userHandler.RegisterHandler)
	r.POST("/login", userHandler.LoginHandler)
	r.POST("/logout", middleware.AuthMiddleware(cfg.JWTSecret), userHandler.LogoutHandler)
	r.POST("/refresh", userHandler.RefreshHandler)

	r.GET("/products", productHandler.GetAllHandler)
	r.GET("/products/:id", productHandler.GetByIDHandler)
	r.POST("/products", productHandler.CreateHandler)
	r.PUT("/products/:id", productHandler.UpdateHandler)
	r.DELETE("/products/:id", productHandler.DeleteHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Запуск сервера в горутине
	go func() {
		zapLogger.Info("HTTP server starting", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	// Ожидание сигнала остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Shutting down server...")

	// Даём время на завершение текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("Server forced to shutdown", zap.Error(err))
	}

	zapLogger.Info("Server exited gracefully")
}
