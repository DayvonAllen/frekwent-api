package handlers

import (
	"fmt"
	"freq/models"
	"freq/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type LoginIpHandler struct {
	LoginIpService services.LoginIpService
}

func (lh *LoginIpHandler) FindAll(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	newLoginQuery := c.Query("new", "false")

	isNew, err := strconv.ParseBool(newLoginQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	ips, err := lh.LoginIpService.FindAll(page, isNew)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": ips})
}

func (lh *LoginIpHandler) FindByIp(c *fiber.Ctx) error {
	c.Accepts("application/json")
	ip := new(models.LoginIP)
	err := c.BodyParser(ip)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	foundCoupon, err := lh.LoginIpService.FindByIp(ip.IpAddress)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": foundCoupon})
}
