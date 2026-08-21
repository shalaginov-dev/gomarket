package handler

import (
	"errors"
	"gomarket/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}
func (h *OrderHandler) Checkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Error(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}
	uid, ok := userID.(int)
	if !ok {
		Error(c, http.StatusBadRequest, errors.New("invalid user ID format"))
		return
	}
	order, err := h.orderService.Checkout(c, uid)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}
	Created(c, order, "Order created successfully")
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Error(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}
	uid, ok := userID.(int)
	if !ok {
		Error(c, http.StatusBadRequest, errors.New("invalid user ID format"))
		return
	}
	order, err := h.orderService.GetAll(c, uid)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(c, order)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		Error(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}
	uid, ok := userID.(int)
	if !ok {
		Error(c, http.StatusBadRequest, errors.New("invalid user ID format"))
		return
	}
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(c, http.StatusBadRequest, err, "Order ID should be a number")
		return
	}
	order, err := h.orderService.GetByID(c, uid, id)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}
	Success(c, order)
}
