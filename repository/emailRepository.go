package repository

import (
	"freq/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmailRepo interface {
	Create(email *models.Email) error
	SendMassEmail(emails *[]string, coupon string) error
	FindAll(string, bool) (*models.EmailList, error)
	FindAllByEmail(string, bool, string) (*models.EmailList, error)
	FindAllByStatus(string, bool, *models.Status) (*models.EmailList, error)
	UpdateEmailStatus(primitive.ObjectID, models.Status) error
}
