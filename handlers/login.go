package handlers

import (
	"fmt"
	"log"
	"os"
	"strings"
	"user_service/db"
	"user_service/security"
	"user_service/types"
	"user_service/util"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Login(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := uuid.New()
		log.Printf("%s | %s\n", util.GetFunctionName(Login), txid.String())
		err_string := fmt.Sprintf("Unauthorized: %s\n", txid.String())
		username, password, has_auth, err := security.GetBasicAuth(c.Get(fiber.HeaderAuthorization), config)
		if has_auth && err == nil {
			stored_password, err := db.GetPassword(username)
			if err != nil {
				return c.Status(fiber.StatusForbidden).SendString(fmt.Sprintf("%s: %s\n", err.Error(), txid.String()))
			}
			err = bcrypt.CompareHashAndPassword([]byte(stored_password), []byte(password))
			if err != nil {
				log.Printf("Invalid password: %s\n", err.Error())
				return c.Status(fiber.StatusUnauthorized).SendString(err_string)
			}
			user_id, err := db.GetUserId(username)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Internal Server Error: %s\n", txid.String()))
			}
			token, err := security.GenerateJWT(txid, user_id, config)
			if err != nil {
				log.Printf("Error generating jwt: %s\n", err.Error())
				return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Internal Server Error: %s\n", txid.String()))
			}
			// Return Authorized
			response := fiber.Map{
				"txid":    txid.String(),
				"user_id": user_id.String(),
				"token":   fmt.Sprintf("Bearer %s", token),
			}
			return c.Status(fiber.StatusOK).JSON(response)
		} else {
			log.Println("Invalid credentials")
			return c.Status(fiber.StatusUnauthorized).SendString(err_string)
		}
	}
}

func Logout(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(Logout), txid.String())

	token := strings.TrimPrefix(c.Get(fiber.HeaderAuthorization), "Bearer ")
	err := security.Logout(token)
	if err != nil {
		log.Println(err.Error())
		err_string := fmt.Sprintf("Unauthorized: %s\n", txid.String())
		return c.Status(fiber.StatusUnauthorized).SendString(err_string)
	}

	user := c.Locals("user").(types.UserDbo)
	response := fiber.Map{
		"txid":    txid.String(),
		"user_id": user.ID.String(),
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func PublicKey(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := uuid.New()
		log.Printf("%s | %s\n", util.GetFunctionName(PublicKey), txid.String())
		if config.App.Host.UseTLS {
			log.Printf("Service is running TLS, no need to send public key.\n")
			return c.Status(fiber.StatusServiceUnavailable).SendString(fmt.Sprintf("Service is running TLS, just log in: %s\n", txid.String()))
		}

		key, err := os.ReadFile(config.App.Host.CertificatePath)
		if err != nil {
			log.Printf("Error reading cert file: %s\n", err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Internal Server Error: %s\n", txid.String()))
		}
		response := fiber.Map{
			"txid":       txid.String(),
			"public_key": string(key),
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func RefreshToken(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := uuid.New()
		log.Printf("%s | %s\n", util.GetFunctionName(RefreshToken), txid.String())
		user_id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			log.Printf("invalid user: %s\n", err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("invalid user: %s\n", txid.String()))
		}
		token, err := security.GenerateJWT(txid, user_id, config)
		if err != nil {
			log.Printf("Error generating jwt: %s\n", err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Internal Server Error: %s\n", txid.String()))
		}
		response := fiber.Map{
			"txid":    txid.String(),
			"user_id": user_id.String(),
			"token":   fmt.Sprintf("Bearer %s", token),
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}
