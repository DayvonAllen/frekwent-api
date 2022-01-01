package repository

import (
	"freq/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MailMemberRepo interface {
	Create(mailMember *models.MailMember) error
	FindAll() (*[]models.MailMember, error)
	DeleteById(id primitive.ObjectID) error
}
