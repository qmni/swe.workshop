package app

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/qmni/swe.workshop/internal/httpapi"
	"github.com/qmni/swe.workshop/internal/middleware"
	"gorm.io/gorm"
)

type errorEnvelope struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func New(db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "swe-workshop",
		ErrorHandler: jsonErrorHandler,
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

func jsonErrorHandler(c *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	message := http.StatusText(status)

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		status = fiberErr.Code
		if fiberErr.Message != "" {
			message = fiberErr.Message
		} else {
			message = http.StatusText(status)
		}
	}

	return c.Status(status).JSON(errorEnvelope{
		Error: apiError{
			Code:    errorCodeForStatus(status),
			Message: message,
		},
	})
}

func errorCodeForStatus(status int) string {
	switch status {
	case fiber.StatusBadRequest:
		return "BAD_REQUEST"
	case fiber.StatusUnauthorized:
		return "UNAUTHORIZED"
	case fiber.StatusForbidden:
		return "FORBIDDEN"
	case fiber.StatusNotFound:
		return "NOT_FOUND"
	case fiber.StatusConflict:
		return "CONFLICT"
	default:
		if status >= 500 {
			return "INTERNAL_ERROR"
		}
		return "HTTP_ERROR"
	}
}
