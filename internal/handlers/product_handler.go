package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinhnglx/go-rest-api-template/internal/services"
)

type ProductHandler struct {
	productService *services.ProductService
}

// Constructor for ProductHandler
func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GET /products
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.productService.GetAllProducts()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}
