package services

import (
	"freq/models"
	"freq/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PurchaseService interface {
	Purchase(purchase *models.Purchase) error
	FindAll(string, bool) (*models.PurchaseList, error)
	FindByPurchaseById(primitive.ObjectID) (*models.Purchase, error)
	FindByPurchaseConfirmationId(string) (*models.Purchase, error)
	UpdateShippedStatus(*models.PurchaseShippedDTO) error
	UpdateDeliveredStatus(*models.PurchaseDeliveredDTO) error
	UpdatePurchaseAddress(*models.PurchaseAddressDTO) error
	UpdateTrackingNumber(*models.PurchaseTrackingDTO) error
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

func (p DefaultPurchaseService) FindAll(page string, newQuery bool) (*models.PurchaseList, error) {
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

func (p DefaultPurchaseService) FindByPurchaseConfirmationId(id string) (*models.Purchase, error) {
	purchase, err := p.repo.FindByPurchaseConfirmationId(id)

	if err != nil {
		return nil, err
	}

	return purchase, nil
}

func (p DefaultPurchaseService) UpdateShippedStatus(shipped *models.PurchaseShippedDTO) error {
	err := p.repo.UpdateShippedStatus(shipped)

	if err != nil {
		return err
	}

	return nil
}

func (p DefaultPurchaseService) UpdateDeliveredStatus(del *models.PurchaseDeliveredDTO) error {
	err := p.repo.UpdateDeliveredStatus(del)

	if err != nil {
		return err
	}

	return nil
}

func (p DefaultPurchaseService) UpdatePurchaseAddress(add *models.PurchaseAddressDTO) error {
	err := p.repo.UpdatePurchaseAddress(add)

	if err != nil {
		return err
	}

	return nil
}

func (p DefaultPurchaseService) UpdateTrackingNumber(trac *models.PurchaseTrackingDTO) error {
	err := p.repo.UpdateTrackingNumber(trac)

	if err != nil {
		return err
	}

	return nil
}

func NewPurchaseService(repository repository.PurchaseRepo) DefaultPurchaseService {
	return DefaultPurchaseService{repository}
}
