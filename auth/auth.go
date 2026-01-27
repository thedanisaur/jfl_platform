package auth

import (
	"crypto/rsa"
	"fmt"
	"log"

	"github.com/thedanisaur/jfl_platform/security"
	"github.com/thedanisaur/jfl_platform/types"
	"github.com/thedanisaur/jfl_platform/util"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AuthenticationMiddleware(config types.Config, public_key *rsa.PublicKey) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := uuid.New()
		log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(AuthenticationMiddleware))

		route := c.Route()
		if route != nil {
			log.Printf("%s | method: %s | path: %s | name: %s", txid.String(), route.Method, route.Path, route.Name)
		}

		claims, err := security.ValidateJWT(c, config, public_key)
		if err != nil {
			log.Printf("%s | Failed to Validate JWT\n", txid.String())
			err_string := fmt.Sprintf("Unauthorized: %s\n", txid.String())
			return c.Status(fiber.StatusInternalServerError).SendString(err_string)
		}
		user_id := claims["user_id"].(uuid.UUID)
		unit_id := claims["unit_id"].(string)
		role_name := claims["role_name"].(string)
		c.Locals("user_id", user_id)
		c.Locals("unit_id", unit_id)
		c.Locals("role_name", role_name)
		c.Locals("transaction_id", txid)
		return c.Next()
	}
}
