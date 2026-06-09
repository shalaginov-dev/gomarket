package handler

import (
	"errors"
	"gomarket/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService *service.CartService
}

type CartRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService,
	}
}

func (h *CartHandler) GetCartHandler(c *gin.Context) {
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
	cart, err := h.cartService.GetItem(c, uid)
	if err != nil {
		Error(c, http.StatusNotFound, err, "Cart item not found")
		return
	}

	Success(c, cart)
}

func (h *CartHandler) AddItemHandler(c *gin.Context) {
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
	var cartReq CartRequest
	// Проверяем тело запроса
	if err := c.ShouldBindJSON(&cartReq); err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}
	err := h.cartService.AddItem(c, uid, cartReq.ProductID, cartReq.Quantity)
	if err != nil {
		Error(c, http.StatusInternalServerError, err, "Cart item could not be added")
		return
	}
	Success(c, nil)
}
func (h *CartHandler) RemoveItemHandler(c *gin.Context) {
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
	productIDSrt := c.Param("product_id")
	productID, err := strconv.Atoi(productIDSrt)
	if err != nil {
		Error(c, http.StatusBadRequest, errors.New("invalid product ID format"))
		return
	}
	err = h.cartService.RemoveItem(c, uid, productID)
	if err != nil {
		Error(c, http.StatusInternalServerError, err, "Cart item could not be removed")
		return
	}
	Success(c, nil)
}
