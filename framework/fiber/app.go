package fiber

import "github.com/gofiber/fiber/v2"

type FiberApp struct {
	*fiber.App
}

func NewFiberApp() *FiberApp {
	app := fiber.New(fiber.Config{
		AppName: "HLAB-Backend",
	})

	return &FiberApp{app}
}
