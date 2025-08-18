package services

import (
	"github.com/vinhnglx/go-rest-api-template/internal/models"
	"github.com/vinhnglx/go-rest-api-template/internal/repositories"
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
