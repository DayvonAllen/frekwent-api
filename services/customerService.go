package services

import (
	"freq/models"
	"freq/repository"
)

type CustomerService interface {
	Create(customer *models.Customer) error
	FindAll(string, bool) (*models.CustomerList, error)
	FindAllByFullName(string, string, string, bool) (*models.CustomerList, error)
}

type DefaultCustomerService struct {
	repo repository.CustomerRepo
}

func (c DefaultCustomerService) Create(customer *models.Customer) error {
	err := c.repo.Create(customer)

	if err != nil {
		return err
	}

	return nil
}

func (c DefaultCustomerService) FindAll(page string, newQuery bool) (*models.CustomerList, error) {
	customers, err := c.repo.FindAll(page, newQuery)

	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (c DefaultCustomerService) FindAllByFullName(firstName string, lastName string,
	page string, newQuery bool) (*models.CustomerList, error) {
	customers, err := c.repo.FindAllByFullName(firstName, lastName, page, newQuery)

	if err != nil {
		return nil, err
	}

	return customers, nil
}

func NewCustomerService(repository repository.CustomerRepo) DefaultCustomerService {
	return DefaultCustomerService{repository}
}
