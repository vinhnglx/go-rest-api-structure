package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vinhnglx/go-rest-api-template/config"
	"github.com/vinhnglx/go-rest-api-template/internal/database"
	"github.com/vinhnglx/go-rest-api-template/internal/handlers"
	"github.com/vinhnglx/go-rest-api-template/internal/middleware"
	"github.com/vinhnglx/go-rest-api-template/internal/repositories"
	"github.com/vinhnglx/go-rest-api-template/internal/services"
)

func main() {
	config.Load()

	database.Connect()

	productRepo := repositories.NewProductRepository(database.DB)

	productService := services.NewProductService(productRepo)

	productHandler := handlers.NewProductHandler(productService)

	router := gin.Default()

	router.Use(middleware.CORS())

	api := router.Group("/api/v1")
	{
		api.GET("/products", productHandler.GetAllProducts)
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
