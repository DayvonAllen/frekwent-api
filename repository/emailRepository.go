package repository

import "freq/models"

type EmailRepo interface {
	Create(email *models.Email) error
	FindAll(string, bool) (*models.EmailList, error)
	FindAllByEmail(string, bool, string) (*models.EmailList, error)
}
