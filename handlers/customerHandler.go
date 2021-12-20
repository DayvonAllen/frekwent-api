package handlers

import (
	"fmt"
	"freq/models"
	"freq/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type CustomerHandler struct {
	CustomerService services.CustomerService
}

func (ch *CustomerHandler) FindAll(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	newCustomerQuery := c.Query("new", "false")

	isNew, err := strconv.ParseBool(newCustomerQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	ips, err := ch.CustomerService.FindAll(page, isNew)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": ips})
}

func (ch *CustomerHandler) FindAllByFullName(c *fiber.Ctx) error {
	c.Accepts("application/json")
	page := c.Query("page", "1")
	newCustomerQuery := c.Query("new", "false")

	customer := new(models.Customer)
	err := c.BodyParser(customer)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	isNew, err := strconv.ParseBool(newCustomerQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	customers, err := ch.CustomerService.FindAllByFullName(customer.FirstName, customer.LastName, page, isNew)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": customers})
}
