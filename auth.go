package main

import (
	"fmt"
	"jfl_platform/security"
	"jfl_platform/types"
	"jfl_platform/util"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AuthorizationMiddleware(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := uuid.New()
		log.Printf("%s | %s\n", util.GetFunctionName(AuthorizationMiddleware), txid.String())

		claims, err := security.ValidateJWT(c, config)
		if err != nil {
			log.Printf("Failed to Validate JWT\n%s\n", err.Error())
			err_string := fmt.Sprintf("Unauthorized: %s\n", txid.String())
			return c.Status(fiber.StatusInternalServerError).SendString(err_string)
		}
		user_id := claims["user_id"].(uuid.UUID)
		c.Locals("user_id", user_id)
		return c.Next()
	}
}
