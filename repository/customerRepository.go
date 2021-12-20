package repository

import "freq/models"

type CustomerRepo interface {
	Create(customer *models.Customer) error
	FindAll(string, bool) (*[]models.Customer, error)
	FindAllByFullName(string, string, string, bool) (*[]models.Customer, error)
}
