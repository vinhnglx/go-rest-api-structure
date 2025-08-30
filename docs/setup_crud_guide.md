### Step 1: Project Structure

```bash
# Create project directory
mkdir my-new-api
cd my-new-api

# Initialize Go module
go mod init github.com/yourusername/my-new-api

# Create directory structure
mkdir -p cmd/{api,worker,migrate,scheduler}
mkdir -p internal/{handlers,services,repositories,models,middleware,database,utils}
mkdir -p {migrations,tests,config,docs,scripts}
mkdir -p tests/{integration,unit,fixtures}
```

### Step 2: Install Dependencies (First Time Only)

**Core Dependencies:**

```bash
# Web framework
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors

# Database & ORM
go get gorm.io/gorm
go get gorm.io/driver/postgres

# Authentication & Security
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt

# Configuration
go get github.com/joho/godotenv

# Validation
go get github.com/go-playground/validator/v10

# Migration dependencies (for cmd/migrate)
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/postgres
go get github.com/golang-migrate/migrate/v4/source/file
go get github.com/lib/pq
```

**Development Tools (Install Globally):**

```bash
# Migration CLI tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Step 3: Environment variables

**Create .env**

```env
# Database
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=myapi_development
DB_PORT=5432
DB_SSLMODE=disable

# Server
PORT=8080
GIN_MODE=debug

```

### Step 4: Create Config & DB setup

**Open config/config.go**

```go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string
	PORT       string
	JWT_SECRET string
}

var AppConfig *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	AppConfig = &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "dbname"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
		PORT:       getEnv("PORT", "8080"),
	}

	log.Println("Configuration loaded successfully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

**Open internal/datbase/connection.go**

```go
package database

import (
	"fmt"
	"log"

	"github.com/vinhnglx/go-rest-api-template/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.AppConfig.DBHost, config.AppConfig.DBUser, config.AppConfig.DBPassword, config.AppConfig.DBName, config.AppConfig.DBPort, config.AppConfig.DBSSLMode)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Println("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

}

func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("Failed to get database connection:", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Println("Failed to close database connection:", err)
	} else {
		log.Println("Database connection closed successfully")
	}
}

```

### Step 5: Create migration script

**Open file cmd/migrate/main.go**

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/yourusername/go-api-template/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/migrate/main.go [up|down|force|version|create]")
	}

	command := os.Args[1]

	// Load configuration
	config.Load()

	// Handle create command separately
	if command == "create" {
		if len(os.Args) < 3 {
			log.Fatal("Usage: go run cmd/migrate/main.go create <migration_name>")
		}
		migrationName := os.Args[2]
		fmt.Printf("Run: migrate create -ext sql -dir migrations -seq %s\n", migrationName)
		return
	}

	// Connect to database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
		config.AppConfig.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create migrate instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create postgres driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}

	// Execute command
	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to run migrations:", err)
		}
		log.Println("Migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to rollback migrations:", err)
		}
		log.Println("Migrations rolled back successfully")
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal("Failed to get version:", err)
		}
		fmt.Printf("Version: %d, Dirty: %t\n", version, dirty)
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Usage: go run cmd/migrate/main.go force <version>")
		}
		// Parse version and force it
		log.Println("Force command - be careful!")
	default:
		log.Fatal("Unknown command. Use: up, down, version, create, or force")
	}
}
```

### Step 6: Create database && migration file

**Create database**

```
createdb myapi_development
```

**Create migration**

```bash
migrate create -ext sql -dir migrations -seq create_users_table
```

**Fill in migration**

```sql
create table if not exists products (
    id serial primary key,
    name varchar(100) not null,
    description text,
    price DOUBLE PRECISION NOT NULL,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);
```

```sql
drop table if exists products;
```

**Run migration**

```bash
go run cmd/migrate/main.go up
```

### Step 7: Create Model

**Product model**

```go
package models

import "time"

type Product struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
}
```

### Step 8: Create Repository

**Product Repository**

```go
package repositories

import (
	"github.com/vinhnglx/go-rest-api-template/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

type ProductRepositoryInterface interface {
	GetAllProducts() ([]models.Product, error)
	GetById(id uint) (*models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(id uint) error
}

// Constructor for ProductRepository
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Instance methods for ProductRepository
func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetById(id uint) (*models.Product, error) {
	var product models.Product

	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

```

### Step 9: Create Service

**Product Service**

```go
package services

import (
	"errors"

	"github.com/vinhnglx/go-rest-api-template/internal/models"
	"github.com/vinhnglx/go-rest-api-template/internal/repositories"
	"gorm.io/gorm"
)

type ProductService struct {
	productRepo repositories.ProductRepositoryInterface
}

// Constructor for ProductService
func NewProductService(repo repositories.ProductRepositoryInterface) *ProductService {
	return &ProductService{
		productRepo: repo,
	}
}

func (s *ProductService) GetAllProducts() ([]models.ProductResponse, error) {
	products, err := s.productRepo.GetAllProducts()

	if err != nil {
		return nil, err
	}

	var response []models.ProductResponse

	for _, product := range products {
		response = append(response, models.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return response, nil
}

func (s *ProductService) GetById(id uint) (*models.ProductResponse, error) {
	product, err := s.productRepo.GetById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	productResponse := &models.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}

	return productResponse, nil
}

func (s *ProductService) Create(req models.CreateProductRequest) (*models.ProductResponse, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	productResponse := &models.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}

	return productResponse, nil
}

func (s *ProductService) Update(id uint, req models.UpdateProductRequest) (*models.ProductResponse, error) {
	existing_product, err := s.productRepo.GetById(id)

	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		existing_product.Name = req.Name
	}

	if req.Description != "" {
		existing_product.Description = req.Description
	}

	if req.Price > 0 {
		existing_product.Price = req.Price
	}

	if err := s.productRepo.Update(existing_product); err != nil {
		return nil, err
	}

	productResponse := &models.ProductResponse{
		ID:          existing_product.ID,
		Name:        existing_product.Name,
		Description: existing_product.Description,
		Price:       existing_product.Price,
	}

	return productResponse, nil
}

func (s *ProductService) Delete(id uint) error {
	return s.productRepo.Delete(id)
}
```

### Step 10: Create Handler

**Product Handler**

```go
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
```

### Step 11: Load the handler

**Open cmd/api/main.go**

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-api-template/config"
	"github.com/yourusername/go-api-template/internal/database"
	"github.com/yourusername/go-api-template/internal/handlers"
	"github.com/yourusername/go-api-template/internal/middleware"
	"github.com/yourusername/go-api-template/internal/repositories"
	"github.com/yourusername/go-api-template/internal/services"
)

func main() {
	// Load configuration
	config.Load()

	// Connect to database (NO AutoMigrate!)
	database.Connect()

	// Initialize repositories
	productRepo := repositories.NewProductRepository(database.DB)

	// Initialize services
	productService := services.NewProductService(productRepo)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)

	// Set up Gin router
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Product routes (RESTful)
		api.GET("/products", productHandler.GetProducts)
		api.POST("/products", productHandler.CreateProduct)
		api.GET("/products/:id", productHandler.GetProduct)
		api.PUT("/products/:id", productHandler.UpdateProduct)
		api.DELETE("/products/:id", productHandler.DeleteProduct)
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", config.AppConfig.Port)
		if err := r.Run(":" + config.AppConfig.Port); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	database.Close()
	log.Println("Server shutdown complete")
}
```

### Step 12: CORS basic middleware

**Open internal/middleware/cors.go**

```go
package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}

	return cors.New(config)
}
```

**Update cmd/api/main.go**

```go
	// Global middleware
	r.Use(middleware.CORS())
```

### Step 13: Run the server

```bash
go run cmd/api/main.go
```
