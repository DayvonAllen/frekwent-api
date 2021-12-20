package services

import (
	"freq/models"
	"freq/repository"
)

type AuthService interface {
	Login(username string, password string, ip string) (*models.User, string, error)
}

type DefaultAuthService struct {
	repo repository.AuthRepo
}

func (a DefaultAuthService) Login(username string, password string, ip string) (*models.User, string, error) {
	u, token, err := a.repo.Login(username, password, ip)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

func NewAuthService(repository repository.AuthRepo) DefaultAuthService {
	return DefaultAuthService{repository}
}
