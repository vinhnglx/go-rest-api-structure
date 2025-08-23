package container

import (
	"github.com/vinhnglx/go-rest-api-template/internal/handlers"
	"github.com/vinhnglx/go-rest-api-template/internal/repositories"
	"github.com/vinhnglx/go-rest-api-template/internal/services"
	"gorm.io/gorm"
)

type Repositories struct {
	Product *repositories.ProductRepository
	User    *repositories.UserRepository
}

type Services struct {
	Product *services.ProductService
	Auth    *services.AuthService
}

type Handlers struct {
	Product *handlers.ProductHandler
	Auth    *handlers.AuthHandler
}

type Container struct {
	DB       *gorm.DB
	Repos    *Repositories
	Services *Services
	Handlers *Handlers
}

func NewContainer(db *gorm.DB) *Container {
	repos := &Repositories{
		Product: repositories.NewProductRepository(db),
		User:    repositories.NewUserRepository(db),
	}

	services := &Services{
		Product: services.NewProductService(repos.Product),
		Auth:    services.NewAuthService(repos.User),
	}

	handlers := &Handlers{
		Product: handlers.NewProductHandler(services.Product),
		Auth:    handlers.NewAuthHandler(services.Auth),
	}

	return &Container{
		DB:       db,
		Repos:    repos,
		Services: services,
		Handlers: handlers,
	}
}
