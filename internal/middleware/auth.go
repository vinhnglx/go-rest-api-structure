package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vinhnglx/go-rest-api-template/internal/services"
)

func RequireAuth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token := tokenParts[1]
		user, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("current_user", user)
		c.Set("userID", user.ID)
		c.Next()
	}
}
