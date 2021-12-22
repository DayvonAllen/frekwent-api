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

	items := purchase.PurchasedItems

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	for _, items := range *items {
		price, err := strconv.ParseInt(items.Price, 10, 16)

		if err != nil {
			panic(err)
		}
		purchase.FinalPrice = int16(price) + purchase.FinalPrice
	}

	purchase.FinalPrice = purchase.FinalPrice + purchase.Tax
	purchase.Id = primitive.NewObjectID()

	purchase.FirstName = strings.ToLower(purchase.FirstName)
	purchase.LastName = strings.ToLower(purchase.LastName)
	purchase.Email = strings.ToLower(purchase.Email)

	err = ph.PurchaseService.Purchase(purchase)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "success", "data": "success"})
}
