package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/qmni/swe.workshop/internal/httpapi"
	"github.com/qmni/swe.workshop/internal/middleware"
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

	players := app.Group("/players", middleware.RequireAuth)
	players.Get("", handler.List)
	players.Get("/:id", handler.Get)
	players.Post("", handler.Create)
	players.Put("/:id", handler.Update)
	players.Delete("/:id", handler.Delete)

	return app
}
