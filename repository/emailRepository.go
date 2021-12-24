package repository

import (
	"freq/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmailRepo interface {
	Create(email *models.Email) error
	FindAll(string, bool) (*models.EmailList, error)
	FindAllByEmail(string, bool, string) (*models.EmailList, error)
	UpdateEmailStatus(primitive.ObjectID, string) error
}
