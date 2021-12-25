package handlers

import (
	"fmt"
	"freq/models"
	"freq/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
)

type PurchaseHandler struct {
	PurchaseService services.PurchaseService
}

func (ph *PurchaseHandler) FindAll(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	newPurchaseQuery := c.Query("new", "false")

	isNew, err := strconv.ParseBool(newPurchaseQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	ips, err := ph.PurchaseService.FindAll(page, isNew)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": ips})
}

func (ph *PurchaseHandler) FindByPurchaseById(c *fiber.Ctx) error {
	c.Accepts("application/json")

	id := c.Params("id")

	monId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	foundPurchase, err := ph.PurchaseService.FindByPurchaseById(monId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": foundPurchase})
}

func (ph *PurchaseHandler) Purchase(c *fiber.Ctx) error {
	c.Accepts("application/json")
	purchase := new(models.Purchase)
	err := c.BodyParser(purchase)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	purchase.Id = primitive.NewObjectID()

	purchase.FirstName = strings.ToLower(purchase.FirstName)
	purchase.LastName = strings.ToLower(purchase.LastName)
	purchase.Email = strings.ToLower(purchase.Email)
	purchase.State = strings.ToLower(purchase.State)
	purchase.City = strings.ToLower(purchase.City)
	purchase.StreetAddress = strings.ToLower(purchase.StreetAddress)
	purchase.OptionalAddress = strings.ToLower(purchase.OptionalAddress)

	err = ph.PurchaseService.Purchase(purchase)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ph *PurchaseHandler) FindByPurchaseConfirmationId(c *fiber.Ctx) error {
	id := c.Params("id")

	product, err := ph.PurchaseService.FindByPurchaseConfirmationId(id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": product})
}

func (ph *PurchaseHandler) CalculateTransactionsByState(c *fiber.Ctx) error {
	state := c.Params("state")

	t, err := ph.PurchaseService.CalculateTransactionsByState(state)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": t})
}

func (ph *PurchaseHandler) UpdateShippedStatus(c *fiber.Ctx) error {
	c.Accepts("application/json")
	purchase := new(models.PurchaseShippedDTO)
	err := c.BodyParser(purchase)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id := c.Params("id")

	monId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	purchase.Id = monId

	err = ph.PurchaseService.UpdateShippedStatus(purchase)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ph *PurchaseHandler) UpdateDeliveredStatus(c *fiber.Ctx) error {
	c.Accepts("application/json")
	purchase := new(models.PurchaseDeliveredDTO)
	err := c.BodyParser(purchase)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id := c.Params("id")

	monId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	purchase.Id = monId

	err = ph.PurchaseService.UpdateDeliveredStatus(purchase)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ph *PurchaseHandler) UpdatePurchaseAddress(c *fiber.Ctx) error {
	c.Accepts("application/json")
	purchase := new(models.PurchaseAddressDTO)
	err := c.BodyParser(purchase)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id := c.Params("id")

	monId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	purchase.Id = monId

	err = ph.PurchaseService.UpdatePurchaseAddress(purchase)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}

func (ph *PurchaseHandler) UpdateTrackingNumber(c *fiber.Ctx) error {
	c.Accepts("application/json")
	purchase := new(models.PurchaseTrackingDTO)
	err := c.BodyParser(purchase)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	id := c.Params("id")

	monId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	purchase.Id = monId

	err = ph.PurchaseService.UpdateTrackingNumber(purchase)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
