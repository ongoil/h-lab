package test

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Test(c *fiber.Ctx) error {
	fmt.Println("test successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "200",
	})
}
