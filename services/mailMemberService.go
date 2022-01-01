package services

import (
	"freq/models"
	"freq/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MailMemberService interface {
	Create(mm *models.MailMember) error
	FindAll() (*[]models.MailMember, error)
	DeleteById(id primitive.ObjectID) error
}

type DefaultMailMemberService struct {
	repo repository.MailMemberRepo
}

func (m DefaultMailMemberService) Create(mm *models.MailMember) error {
	err := m.repo.Create(mm)

	if err != nil {
		return err
	}

	return nil
}

func (m DefaultMailMemberService) FindAll() (*[]models.MailMember, error) {
	members, err := m.repo.FindAll()

	if err != nil {
		return nil, err
	}

	return members, nil
}

func (m DefaultMailMemberService) DeleteById(id primitive.ObjectID) error {
	err := m.repo.DeleteById(id)

	if err != nil {
		return err
	}

	return nil
}

func NewMailMemberService(repository repository.MailMemberRepo) DefaultMailMemberService {
	return DefaultMailMemberService{repository}
}
