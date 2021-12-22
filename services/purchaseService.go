package services

import (
	"freq/models"
	"freq/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PurchaseService interface {
	Purchase(purchase *models.Purchase) error
	FindAll(string, bool) (*[]models.Purchase, error)
	FindByPurchaseById(primitive.ObjectID) (*models.Purchase, error)
}

type DefaultPurchaseService struct {
	repo repository.PurchaseRepo
}

func (p DefaultPurchaseService) Purchase(purchase *models.Purchase) error {
	err := p.repo.Purchase(purchase)

	if err != nil {
		return err
	}

	return nil
}

func (p DefaultPurchaseService) FindAll(page string, newQuery bool) (*[]models.Purchase, error) {
	purchases, err := p.repo.FindAll(page, newQuery)

	if err != nil {
		return nil, err
	}

	return purchases, nil
}

func (p DefaultPurchaseService) FindByPurchaseById(id primitive.ObjectID) (*models.Purchase, error) {
	purchase, err := p.repo.FindByPurchaseById(id)

	if err != nil {
		return nil, err
	}

	return purchase, nil
}

func NewPurchaseService(repository repository.PurchaseRepo) DefaultPurchaseService {
	return DefaultPurchaseService{repository}
}
