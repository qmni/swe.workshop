package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/qmni/swe.workshop/internal/httpapi"
	"gorm.io/gorm"
)

func New(db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "swe-workshop",
	})

	validator := validator.New()
	handler := httpapi.NewPlayerHandler(db, validator)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	app.Get("/players", handler.List)
	app.Get("/players/:id", handler.Get)
	app.Post("/players", handler.Create)

	return app
}
