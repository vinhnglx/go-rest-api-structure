package services

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vinhnglx/go-rest-api-template/internal/models"
	"github.com/vinhnglx/go-rest-api-template/internal/services"
	"github.com/vinhnglx/go-rest-api-template/tests/mocks"
)

type ProductServiceTestSuite struct {
	suite.Suite
	productService  *services.ProductService
	mockProductRepo *mocks.MockProductRepository
}

func (suite *ProductServiceTestSuite) SetupTest() {
	suite.mockProductRepo = new(mocks.MockProductRepository)
	suite.productService = services.NewProductService(suite.mockProductRepo)
}

func (suite *ProductServiceTestSuite) TestGetAllProducts() {
	suite.mockProductRepo.On("GetAllProducts").Return([]models.Product{
		{ID: 1, Name: "Product 1", Description: "Description 1", Price: 100},
		{ID: 2, Name: "Product 2", Description: "Description 2", Price: 200},
	}, nil)

	products, err := suite.productService.GetAllProducts()
	suite.NoError(err)
	suite.Len(products, 2)
	suite.Equal("Product 1", products[0].Name)
	suite.Equal("Product 2", products[1].Name)
}

func TestProductServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}
