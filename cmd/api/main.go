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
		api.GET("/products", container.Handlers.Product.GetAllProducts)
		api.GET("/products/:id", container.Handlers.Product.GetById)
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
