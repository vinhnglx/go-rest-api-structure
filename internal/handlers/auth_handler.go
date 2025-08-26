package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/vinhnglx/go-rest-api-template/internal/models"
	"github.com/vinhnglx/go-rest-api-template/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user)
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	authResponse, err := h.authService.Login(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, authResponse)
}

// GET /auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	user, exists := c.Get("current_user")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	current_user := user.(*models.User)

	response := models.UserResponse{
		Name:  current_user.Name,
		Email: current_user.Email,
	}

	c.JSON(200, response)
}
