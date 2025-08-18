package services

import (
	"errors"

	"github.com/vinhnglx/go-rest-api-template/internal/models"
	"github.com/vinhnglx/go-rest-api-template/internal/repositories"
	"gorm.io/gorm"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

// Constructor for ProductService
func NewProductService(productRepo *repositories.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
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

	return product, nil
}
