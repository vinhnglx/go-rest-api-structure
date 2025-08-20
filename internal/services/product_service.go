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
