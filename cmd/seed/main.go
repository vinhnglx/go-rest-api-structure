package main

import (
	"log"

	"github.com/vinhnglx/go-rest-api-template/config"
	"github.com/vinhnglx/go-rest-api-template/internal/database"
	"github.com/vinhnglx/go-rest-api-template/internal/models"
)

func main() {
	config.Load()

	database.Connect()

	products := []models.Product{
		{
			Name:        "Product 1",
			Description: "Description for Product 1",
			Price:       10.99,
		},
		{
			Name:        "Product 2",
			Description: "Description for Product 2",
			Price:       12.99,
		},
		{
			Name:        "Product 3",
			Description: "Description for Product 3",
			Price:       15.99,
		},
	}

	for _, product := range products {
		if err := database.DB.Create(&product).Error; err != nil {
			log.Fatalf("Failed to seed product %s: %v", product.Name, err)
		}
	}

	log.Println("Database seeded successfully")
}
