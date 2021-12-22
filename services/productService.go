package services

import (
	"freq/models"
	"freq/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService interface {
	Create(product *models.Product) error
	FindAll(string, bool) (*[]models.Product, error)
	FindByProductId(primitive.ObjectID) (*models.Product, error)
	UpdateById(product *models.Product) error
	DeleteById(primitive.ObjectID) error
}

type DefaultProductService struct {
	repo repository.ProductRepo
}

func (p DefaultProductService) Create(product *models.Product) error {
	err := p.repo.Create(product)

	if err != nil {
		return err
	}

	return nil
}

func (p DefaultProductService) FindAll(page string, newQuery bool) (*[]models.Product, error) {
	products, err := p.repo.FindAll(page, newQuery)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p DefaultProductService) FindByProductId(id primitive.ObjectID) (*models.Product, error) {
	product, err := p.repo.FindByProductId(id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p DefaultProductService) UpdateById(product *models.Product) error {
	err := p.repo.UpdateById(product)

	if err != nil {
		return err
	}

	return nil
}

func (p DefaultProductService) DeleteById(id primitive.ObjectID) error {
	err := p.DeleteById(id)

	if err != nil {
		return err
	}

	return nil
}

func NewProductService(repository repository.ProductRepo) DefaultProductService {
	return DefaultProductService{repository}
}
