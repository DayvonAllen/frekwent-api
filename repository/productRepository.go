package repository

import (
	"freq/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRepo interface {
	Create(product *models.Product) error
	FindAll(string, bool) (*[]models.Product, error)
	FindByProductId(primitive.ObjectID) (*models.Product, error)
	UpdateName(string, primitive.ObjectID) error
	UpdateQuantity(uint16, primitive.ObjectID) error
	UpdatePrice(string, primitive.ObjectID) error
	UpdateDescription(string, primitive.ObjectID) error
	UpdateIngredients(*[]string, primitive.ObjectID) error
	DeleteById(primitive.ObjectID) error
}
