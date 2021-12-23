package services

import (
	"freq/models"
	"freq/repository"
)

type EmailService interface {
	FindAll(string, bool) (*models.EmailList, error)
	FindAllByEmail(string, bool, string) (*models.EmailList, error)
}

func (e DefaultEmailService) FindAll(page string, newQuery bool) (*models.EmailList, error) {
	emails, err := e.repo.FindAll(page, newQuery)

	if err != nil {
		return nil, err
	}

	return emails, nil
}

func (e DefaultEmailService) FindAllByEmail(page string, newQuery bool, email string) (*models.EmailList, error) {
	emails, err := e.repo.FindAllByEmail(page, newQuery, email)

	if err != nil {
		return nil, err
	}

	return emails, nil
}

type DefaultEmailService struct {
	repo repository.EmailRepo
}

func NewEmailService(repository repository.EmailRepo) DefaultEmailService {
	return DefaultEmailService{repository}
}
