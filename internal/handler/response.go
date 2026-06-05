package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response — единая структура ответа для всего API
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success — успешный ответ (200)
func Success(c *gin.Context, data interface{}, message ...string) {
	msg := "success"
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

// Created — ресурс успешно создан (201)
func Created(c *gin.Context, data interface{}, message ...string) {
	msg := "created successfully"
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

// Error — универсальная ошибка
func Error(c *gin.Context, statusCode int, err error, message ...string) {
	msg := "something went wrong"
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(statusCode, Response{
		Success: false,
		Message: msg,
		Error:   err.Error(),
	})
}
