package main

import (
	"context"
	"fmt"
	"log"

	"gomarket/internal/config"
	"gomarket/internal/db"

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

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DBDSN)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer pool.Close()
	fmt.Println("Database connected")

	if err := db.RunMigrations(cfg.DBDSN, "./migrations"); err != nil {
		log.Fatal("Failed to run migrations", err)
	}
	fmt.Println("Migrations applied")

	r := gin.Default()
	r.GET("/health", healthHandler)

	if err := r.Run(":" + cfg.Port); err != nil {
		fmt.Println("Server failed:", err)
	}
}
