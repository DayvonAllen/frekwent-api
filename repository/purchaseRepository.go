package repository

import (
	"freq/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PurchaseRepo interface {
	Purchase(purchase *models.Purchase) error
	FindAll(string, bool) (*[]models.Purchase, error)
	FindByPurchaseById(primitive.ObjectID) (*models.Purchase, error)
}
