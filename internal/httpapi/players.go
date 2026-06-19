package httpapi

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/qmni/swe.workshop/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PlayerHandler struct {
	db       *gorm.DB
	validate *validator.Validate
}

type createPlayerRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=60"`
	Email       string `json:"email" validate:"required,email,max=120"`
	Level       int    `json:"level" validate:"omitempty,gte=1,lte=100"`
	Experience  int    `json:"experience" validate:"omitempty,gte=0"`
	PlayerClass string `json:"playerClass" validate:"required,oneof=WARRIOR MAGE ROGUE PRIEST HUNTER"`
	GuildID     *uint  `json:"guildId" validate:"omitempty,gte=1"`
}

type updatePlayerRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=60"`
	Email       string `json:"email" validate:"required,email,max=120"`
	Level       int    `json:"level" validate:"required,gte=1,lte=100"`
	Experience  int    `json:"experience" validate:"gte=0"`
	PlayerClass string `json:"playerClass" validate:"required,oneof=WARRIOR MAGE ROGUE PRIEST HUNTER"`
	Status      string `json:"status" validate:"required,oneof=ACTIVE BANNED DELETED"`
	GuildID     *uint  `json:"guildId" validate:"omitempty,gte=1"`
}

type listPlayersResponse struct {
	Items      []model.Player `json:"items"`
	Pagination paginationMeta `json:"pagination"`
}

type paginationMeta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type errorEnvelope struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func NewPlayerHandler(db *gorm.DB, validate *validator.Validate) PlayerHandler {
	return PlayerHandler{db: db, validate: validate}
}

func (h PlayerHandler) List(c *fiber.Ctx) error {
	page, err := parsePositiveIntQuery(c, "page", 1)
	if err != nil {
		return respondValidationError(c, []string{err.Error()})
	}

	limit, err := parsePositiveIntQuery(c, "limit", 20)
	if err != nil {
		return respondValidationError(c, []string{err.Error()})
	}
	if limit > 100 {
		return respondValidationError(c, []string{"limit must be <= 100"})
	}

	playerClass := c.Query("playerClass")
	if playerClass != "" {
		if err := h.validate.Var(playerClass, "oneof=WARRIOR MAGE ROGUE PRIEST HUNTER"); err != nil {
			return respondValidationError(c, []string{"playerClass must be one of WARRIOR, MAGE, ROGUE, PRIEST, HUNTER"})
		}
	}

	query := h.db.Model(&model.Player{})
	if playerClass != "" {
		query = query.Where("\"playerClass\" = ?", playerClass)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "players could not be loaded")
	}

	var players []model.Player
	offset := (page - 1) * limit
	if err := query.Order("id asc").Offset(offset).Limit(limit).Find(&players).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "players could not be loaded")
	}

	return c.JSON(listPlayersResponse{
		Items: players,
		Pagination: paginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	})
}

func (h PlayerHandler) Get(c *fiber.Ctx) error {
	var player model.Player
	if err := h.db.First(&player, c.Params("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "player not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "player could not be loaded")
	}

	return c.JSON(player)
}

func (h PlayerHandler) Create(c *fiber.Ctx) error {
	var req createPlayerRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}
	if err := h.validate.Struct(req); err != nil {
		return respondValidationError(c, validationErrors(err))
	}

	player := model.Player{
		Username:    req.Username,
		Email:       req.Email,
		Level:       defaultLevel(req.Level),
		Experience:  req.Experience,
		PlayerClass: model.PlayerClass(req.PlayerClass),
		Status:      model.PlayerStatusActive,
		GuildID:     req.GuildID,
	}
	result := h.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&player)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "player could not be created")
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusConflict, "player username or email already exists")
	}

	return c.Status(fiber.StatusCreated).JSON(player)
}

func (h PlayerHandler) Update(c *fiber.Ctx) error {
	var req updatePlayerRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}
	if err := h.validate.Struct(req); err != nil {
		return respondValidationError(c, validationErrors(err))
	}

	var player model.Player
	if err := h.db.First(&player, c.Params("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "player not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "player could not be loaded")
	}

	player.Username = req.Username
	player.Email = req.Email
	player.Level = req.Level
	player.Experience = req.Experience
	player.PlayerClass = model.PlayerClass(req.PlayerClass)
	player.Status = model.PlayerStatus(req.Status)
	player.GuildID = req.GuildID
	player.Version++

	if err := h.db.Save(&player).Error; err != nil {
		if isUniqueViolation(err) {
			return fiber.NewError(fiber.StatusConflict, "player username or email already exists")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "player could not be updated")
	}

	return c.JSON(player)
}

func (h PlayerHandler) Delete(c *fiber.Ctx) error {
	result := h.db.Delete(&model.Player{}, c.Params("id"))
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "player could not be deleted")
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "player not found")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func defaultLevel(level int) int {
	if level == 0 {
		return 1
	}
	return level
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
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

func parsePositiveIntQuery(c *fiber.Ctx, key string, fallback int) (int, error) {
	raw := c.Query(key)
	if raw == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return 0, errors.New(key + " must be a positive integer")
	}

	return value, nil
}

func respondValidationError(c *fiber.Ctx, details []string) error {
	return c.Status(fiber.StatusBadRequest).JSON(errorEnvelope{
		Error: apiError{
			Code:    "VALIDATION_ERROR",
			Message: "validation failed",
			Details: details,
		},
	})
}
