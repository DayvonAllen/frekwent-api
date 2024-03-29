package middleware

import (
	"fmt"
	"freq/util"
	"github.com/gofiber/fiber/v2"
)

func IsLoggedIn(c *fiber.Ctx) error {
	token := c.Cookies("auth")

	var auth util.Authentication
	u, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("Unauthorized user")})
	}

	c.Locals("username", u.Username)
	c.Locals("id", u.Id)

	err = c.Next()

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("Unauthorized user")})
	}

	return nil
}
