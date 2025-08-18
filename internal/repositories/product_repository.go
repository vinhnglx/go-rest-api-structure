package repositories

import (
	"github.com/vinhnglx/go-rest-api-template/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

// Constructor for ProductRepository
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetById(id uint) (*models.ProductResponse, error) {
	var product models.Product

	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	response := &models.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}

	return response, nil
}
