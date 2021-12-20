package services

import (
	"freq/models"
	"freq/repository"
)

type CouponService interface {
	Create(coupon *models.Coupon) error
	FindAll(string, bool) (*[]models.Coupon, error)
	FindByCode(string) (*models.Coupon, error)
	DeleteByCode(string) error
}

type DefaultCouponService struct {
	repo repository.CouponRepo
}

func (c DefaultCouponService) Create(coupon *models.Coupon) error {
	err := c.Create(coupon)

	if err != nil {
		return err
	}

	return nil
}

func (c DefaultCouponService) FindAll(page string, newQuery bool) (*[]models.Coupon, error) {
	coupons, err := c.FindAll(page, newQuery)

	if err != nil {
		return nil, err
	}

	return coupons, nil
}

func (c DefaultCouponService) FindByCode(code string) (*models.Coupon, error) {
	coupon, err := c.FindByCode(code)

	if err != nil {
		return nil, err
	}

	return coupon, nil
}

func (c DefaultCouponService) DeleteByCode(code string) error {
	err := c.DeleteByCode(code)

	if err != nil {
		return err
	}

	return nil
}

func NewCouponService(repository repository.CouponRepo) DefaultCouponService {
	return DefaultCouponService{repository}
}
