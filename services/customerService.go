package services

import (
	"freq/models"
	"freq/repository"
)

type CustomerService interface {
	Create(customer *models.Customer) error
	FindAll(string, bool) (*[]models.Customer, error)
	FindAllByFullName(string, string, string, bool) (*[]models.Customer, error)
}

type DefaultCustomerService struct {
	repo repository.CustomerRepo
}

func (c DefaultCustomerService) Create(customer *models.Customer) error {
	err := c.Create(customer)

	if err != nil {
		return err
	}

	return nil
}

func (c DefaultCustomerService) FindAll(page string, newQuery bool) (*[]models.Customer, error) {
	ips, err := c.FindAll(page, newQuery)

	if err != nil {
		return nil, err
	}

	return ips, nil
}

func (c DefaultCustomerService) FindAllByFullName(firstName string, lastName string,
	page string, newQuery bool) (*[]models.Customer, error) {
	ips, err := c.FindAllByFullName(firstName, lastName, page, newQuery)

	if err != nil {
		return nil, err
	}

	return ips, nil
}

func NewCustomerService(repository repository.CustomerRepo) DefaultCustomerService {
	return DefaultCustomerService{repository}
}
