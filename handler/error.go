package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TODO [drd] make sure the services implement this for the standard
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "internal server error"

	if fiber_error, ok := err.(*fiber.Error); ok {
		code = fiber_error.Code
		msg = fiber_error.Message
	}

	response := fiber.Map{
		"error": msg,
	}

	transaction_id, ok := c.Locals("transaction_id").(uuid.UUID)
	if ok {
		response["transaction_id"] = transaction_id.String()
	}

	return c.Status(code).JSON(response)
}
