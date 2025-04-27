package main

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type Handler struct {
	db  db
	log *zap.Logger
}

func NewHandler(db db, logger *zap.Logger) *Handler {
	return &Handler{db: db, log: logger}
}

func (h *Handler) CreateTimer(c *fiber.Ctx) error {
	var req CreateTimerRequest
	h.log.Debug("starting timer", zap.String("key", "val"))
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	timer := Timer{
		TimeInterval: req.TimeInterval,
		Pauses:       req.Pauses,
	}

	id, err := h.db.CreateTimer(timer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create timer",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}

func (h *Handler) GetTimer(c *fiber.Ctx) error {
	idStr := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid timer ID",
		})
	}

	timer, err := h.db.GetTimer(objectID.Hex())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Timer not found",
		})
	}

	return c.JSON(timer)
}

func (h *Handler) UpdateTimer(c *fiber.Ctx) error {
	id := c.Params("id")

	var req CreateTimerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	timer := Timer{
		TimeInterval: req.TimeInterval,
		Pauses:       req.Pauses,
	}

	if err := h.db.UpdateTimer(id, timer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update timer",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) DeleteTimer(c *fiber.Ctx) error {
	idStr := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid timer ID",
		})
	}

	if err := h.db.DeleteTimer(objectID.Hex()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete timer",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) ListTimers(c *fiber.Ctx) error {
	h.log.Info("called")
	timers, err := h.db.ListTimers()
	h.log.Info("got timers", zap.Any("timers", timers))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list timers",
		})
	}

	return c.JSON(timers)
}
