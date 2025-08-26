package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vinhnglx/go-rest-api-template/config"
	"github.com/vinhnglx/go-rest-api-template/internal/container"
	"github.com/vinhnglx/go-rest-api-template/internal/database"
	"github.com/vinhnglx/go-rest-api-template/internal/middleware"
)

func main() {
	config.Load()

	database.Connect()

	container := container.NewContainer(database.DB)

	router := gin.Default()

	router.Use(middleware.CORS())

	api := router.Group("/api/v1")
	{
		productProtected := api.Group("/products").Use(middleware.RequireAuth(container.Services.Auth))

		productProtected.GET("/", container.Handlers.Product.GetAllProducts)
		productProtected.GET("/:id", container.Handlers.Product.GetById)
		productProtected.POST("/", container.Handlers.Product.Create)
		productProtected.PUT("/:id", container.Handlers.Product.Update)
		productProtected.DELETE("/:id", container.Handlers.Product.Delete)

		auth := api.Group("/auth")
		{
			auth.POST("/login", container.Handlers.Auth.Login)
			auth.POST("/register", container.Handlers.Auth.Register)
		}

		authProtected := auth.Group("/").Use(middleware.RequireAuth(container.Services.Auth))
		{
			authProtected.GET("/me", container.Handlers.Auth.Me)
		}
	}

	go func() {
		log.Printf("Server is starting on port %s", config.AppConfig.PORT)

		if err := router.Run(":" + config.AppConfig.PORT); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan struct{})
	<-quit

	log.Println("Server stopped gracefully")
	database.Close()
	log.Println("Database connection closed")
}
