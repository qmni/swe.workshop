package httpapi

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/qmni/swe.workshop/internal/model"
	"gorm.io/gorm"
)

type ProductHandler struct {
	db       *gorm.DB
	validate *validator.Validate
}

type createProductRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=120"`
	Description string `json:"description" validate:"max=500"`
	PriceCents  int    `json:"priceCents" validate:"required,gte=1,lte=10000000"`
}

func NewProductHandler(db *gorm.DB, validate *validator.Validate) ProductHandler {
	return ProductHandler{db: db, validate: validate}
}

func (h ProductHandler) List(c *fiber.Ctx) error {
	var products []model.Product
	if err := h.db.Order("id asc").Find(&products).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "products could not be loaded")
	}

	return c.JSON(products)
}

func (h ProductHandler) Get(c *fiber.Ctx) error {
	var product model.Product
	if err := h.db.First(&product, c.Params("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "product not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "product could not be loaded")
	}

	return c.JSON(product)
}

func (h ProductHandler) Create(c *fiber.Ctx) error {
	var req createProductRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation failed",
			"details": validationErrors(err),
		})
	}

	product := model.Product{
		Name:        req.Name,
		Description: req.Description,
		PriceCents:  req.PriceCents,
	}
	if err := h.db.Create(&product).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "product could not be created")
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func validationErrors(err error) []string {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return []string{err.Error()}
	}

	result := make([]string, 0, len(validationErrs))
	for _, fieldErr := range validationErrs {
		result = append(result, fieldErr.Field()+" failed "+fieldErr.Tag())
	}
	return result
}
