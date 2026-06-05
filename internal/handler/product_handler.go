package handler

import (
	"gomarket/internal/domain"
	"gomarket/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) CreateHandler(c *gin.Context) {
	var p domain.Product
	// Проверяем тело запроса
	if err := c.ShouldBindJSON(&p); err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}
	// Создаем продукт
	product, err := h.productService.Create(c, &p)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}
	Created(c, product, "Product created successfully")
}

func (h *ProductHandler) GetAllHandler(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
		if limit > 50 {
			limit = 50 // защита от слишком большого лимита
		}
	}
	products, err := h.productService.GetAll(c, limit)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(c, products)
}

func (h *ProductHandler) GetByIDHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(c, http.StatusBadRequest, err, "Product ID should be a number")
		return
	}

	product, err := h.productService.GetByID(c, id)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(c, product)
}

func (h *ProductHandler) UpdateHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(c, http.StatusBadRequest, err, "Product ID should be a number")
		return
	}
	var p domain.Product
	// Проверяем тело запроса
	if err := c.ShouldBindJSON(&p); err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}
	// Создаем продукт
	product, err := h.productService.Update(c, &p, id)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}
	Success(c, product, "Product updated successfully")
}

func (h *ProductHandler) DeleteHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(c, http.StatusBadRequest, err, "Product ID should be a number")
		return
	}

	if err := h.productService.Delete(c, id); err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(c, "", "Product deleted successfully")
}
