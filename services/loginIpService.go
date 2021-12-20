package services

import (
	"freq/models"
	"freq/repository"
)

type LoginIpService interface {
	Create(ip *models.LoginIP) error
	FindAll(string, bool) (*[]models.LoginIP, error)
	FindByIp(string) (*models.LoginIP, error)
	UpdateIp(ip *models.LoginIP) error
}

type DefaultLoginIpService struct {
	repo repository.LoginIpRepo
}

func (l DefaultLoginIpService) Create(ip *models.LoginIP) error {
	err := l.Create(ip)

	if err != nil {
		return err
	}

	return nil
}

func (l DefaultLoginIpService) FindAll(page string, newQuery bool) (*[]models.LoginIP, error) {
	ips, err := l.FindAll(page, newQuery)

	if err != nil {
		return nil, err
	}

	return ips, nil
}

func (l DefaultLoginIpService) FindByIp(ip string) (*models.LoginIP, error) {
	i, err := l.FindByIp(ip)

	if err != nil {
		return nil, err
	}

	return i, nil
}

func (l DefaultLoginIpService) UpdateIp(ip *models.LoginIP) error {
	err := l.UpdateIp(ip)

	if err != nil {
		return err
	}

	return nil
}

func NewLoginIpService(repository repository.LoginIpRepo) DefaultLoginIpService {
	return DefaultLoginIpService{repository}
}
