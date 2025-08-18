package container

import (
	"github.com/vinhnglx/go-rest-api-template/internal/handlers"
	"github.com/vinhnglx/go-rest-api-template/internal/repositories"
	"github.com/vinhnglx/go-rest-api-template/internal/services"
	"gorm.io/gorm"
)

type Repositories struct {
	Product *repositories.ProductRepository
}

type Services struct {
	Product *services.ProductService
}

type Handlers struct {
	Product *handlers.ProductHandler
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
	}

	services := &Services{
		Product: services.NewProductService(repos.Product),
	}

	handlers := &Handlers{
		Product: handlers.NewProductHandler(services.Product),
	}

	return &Container{
		DB:       db,
		Repos:    repos,
		Services: services,
		Handlers: handlers,
	}
}
