package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vinhnglx/go-rest-api-template/internal/models"
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

// GET /products/:id
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

// POST /products
func (h *ProductHandler) Create(c *gin.Context) {
	var req models.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	product, err := h.productService.Create(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": product})
}

// PUT /products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	var req models.UpdateProductRequest

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
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

	updated_product, err := h.productService.Update(uint(id), models.UpdateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updated_product})
}

// DELETE /products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.productService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
