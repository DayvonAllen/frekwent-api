package handlers

import (
	"fmt"
	"freq/models"
	"freq/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type CouponHandler struct {
	CouponService services.CouponService
}

func (lh *CouponHandler) Create(c *fiber.Ctx) error {
	c.Accepts("application/json")
	coupon := new(models.Coupon)
	err := c.BodyParser(coupon)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = lh.CouponService.Create(coupon)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (lh *CouponHandler) FindAll(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	newCouponQuery := c.Query("new", "false")

	isNew, err := strconv.ParseBool(newCouponQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	coupons, err := lh.CouponService.FindAll(page, isNew)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": coupons})
}

func (lh *CouponHandler) FindByCode(c *fiber.Ctx) error {
	c.Accepts("application/json")
	coupon := new(models.Coupon)
	err := c.BodyParser(coupon)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	foundCoupon, err := lh.CouponService.FindByCode(coupon.Code)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": foundCoupon})
}

func (lh *CouponHandler) DeleteByCode(c *fiber.Ctx) error {
	c.Accepts("application/json")
	coupon := new(models.Coupon)
	err := c.BodyParser(coupon)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = lh.CouponService.DeleteByCode(coupon.Code)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
