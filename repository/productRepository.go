package repository

import (
	"freq/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRepo interface {
	Create(product *models.Product) error
	FindAll(string, bool) (*[]models.Product, error)
	FindByProductId(primitive.ObjectID) (*models.Product, error)
	UpdateById(product *models.Product) error
	DeleteById(primitive.ObjectID) error
}
