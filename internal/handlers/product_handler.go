package handlers

import (
	"net/http"
	"strconv"

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

func (h *ProductHandler) GetById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	product, err := h.productService.GetById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}
