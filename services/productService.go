package services

import (
	"freq/models"
	"freq/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService interface {
	Create(product *models.Product) error
	FindAll(string, bool) (*models.ProductList, error)
	FindByProductId(primitive.ObjectID) (*models.Product, error)
	UpdateName(string, primitive.ObjectID) (*models.Product, error)
	UpdateQuantity(uint16, primitive.ObjectID) (*models.Product, error)
	UpdatePrice(string, primitive.ObjectID) (*models.Product, error)
	UpdateDescription(string, primitive.ObjectID) (*models.Product, error)
	UpdateIngredients(*[]string, primitive.ObjectID) (*models.Product, error)
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

func (p DefaultProductService) FindAll(page string, newQuery bool) (*models.ProductList, error) {
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

func (p DefaultProductService) UpdateName(name string, id primitive.ObjectID) (*models.Product, error) {
	product, err := p.repo.UpdateName(name, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p DefaultProductService) UpdateQuantity(quan uint16, id primitive.ObjectID) (*models.Product, error) {
	product, err := p.repo.UpdateQuantity(quan, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p DefaultProductService) UpdatePrice(price string, id primitive.ObjectID) (*models.Product, error) {
	product, err := p.repo.UpdatePrice(price, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p DefaultProductService) UpdateDescription(desc string, id primitive.ObjectID) (*models.Product, error) {
	product, err := p.repo.UpdateDescription(desc, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p DefaultProductService) UpdateIngredients(ingredients *[]string, id primitive.ObjectID) (*models.Product, error) {
	product, err := p.repo.UpdateIngredients(ingredients, id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p DefaultProductService) DeleteById(id primitive.ObjectID) error {
	err := p.repo.DeleteById(id)

	if err != nil {
		return err
	}

	return nil
}

func NewProductService(repository repository.ProductRepo) DefaultProductService {
	return DefaultProductService{repository}
}
