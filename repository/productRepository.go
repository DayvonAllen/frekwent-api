package repository

import (
	"freq/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRepo interface {
	Create(product *models.Product) error
	FindAll(string, bool) (*models.ProductList, error)
	FindByProductId(primitive.ObjectID) (*models.Product, error)
	FindByProductName(string) (*models.Product, error)
	FindAllByCategory(string, string, bool) (*models.ProductList, error)
	UpdatePurchaseCount(string) error
	UpdateName(string, primitive.ObjectID) (*models.Product, error)
	UpdateQuantity(uint16, primitive.ObjectID) (*models.Product, error)
	UpdatePrice(string, primitive.ObjectID) (*models.Product, error)
	UpdateDescription(string, primitive.ObjectID) (*models.Product, error)
	UpdateIngredients(*[]string, primitive.ObjectID) (*models.Product, error)
	UpdateCategory(string, primitive.ObjectID) (*models.Product, error)
	DeleteById(primitive.ObjectID) error
}
