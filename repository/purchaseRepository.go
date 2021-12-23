package repository

import (
	"freq/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PurchaseRepo interface {
	Purchase(purchase *models.Purchase) error
	FindAll(string, bool) (*models.PurchaseList, error)
	FindByPurchaseById(primitive.ObjectID) (*models.Purchase, error)
	FindByPurchaseConfirmationId(string) (*models.Purchase, error)
	UpdateShippedStatus(*models.PurchaseShippedDTO) error
	UpdateDeliveredStatus(*models.PurchaseDeliveredDTO) error
	UpdatePurchaseAddress(*models.PurchaseAddressDTO) error
	UpdateTrackingNumber(*models.PurchaseTrackingDTO) error
}
