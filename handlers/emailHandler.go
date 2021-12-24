package handlers

import (
	"fmt"
	"freq/helper"
	"freq/models"
	"freq/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

type EmailHandler struct {
	EmailService services.EmailService
}

func (eh *EmailHandler) FindAll(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	newEmailQuery := c.Query("new", "false")

	isNew, err := strconv.ParseBool(newEmailQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	emails, err := eh.EmailService.FindAll(page, isNew)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": emails})
}

func (eh *EmailHandler) FindAllByEmail(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	newEmailQuery := c.Query("new", "false")
	email := c.Params("email")

	isNew, err := strconv.ParseBool(newEmailQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	if !helper.IsEmail(email) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("invalid email")})
	}

	emails, err := eh.EmailService.FindAllByEmail(page, isNew, strings.ToLower(email))

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": emails})
}

func (eh *EmailHandler) SendEmail(c *fiber.Ctx) error {
	emailType := c.Params("emailType")
	c.Accepts("application/json")
	email := new(models.EmailDto)
	err := c.BodyParser(email)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	createdEmail := helper.CreateEmail(new(models.Email), email, emailType)

	err = eh.EmailService.Create(createdEmail)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "created"})
}
