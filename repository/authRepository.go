package repository

import "freq/models"

type AuthRepo interface {
	Login(username string, password string, ip string) (*models.User, string, error)
}
