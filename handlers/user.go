package handlers

import (
	"fmt"
	"log"
	"user_service/db"
	"user_service/types"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("CreateUser | %s\n", txid.String())
	// DBO here because we're creating the object for the first time
	var user types.UserDbo
	err := c.BodyParser(&user)
	if err != nil {
		log.Printf("Failed to parse user data\n%s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Failed to parse user data: %s\n", txid.String()))
	}
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), 12)
	if err != nil {
		log.Printf("Failed to hash password\n%s\n", err.Error())
		err_string := fmt.Sprintf("Internal Server Error: %s\n", txid.String())
		return c.Status(fiber.StatusInternalServerError).SendString(err_string)
	}
	id, err := db.InsertUser(string(hashed_password), user)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(id)
}

func GetUser(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("GetUser | %s\n", txid.String())
	user_id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).SendString("invalid user")
	}
	user, err := db.GetUser(user_id)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func GetUsers(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("GetUsers | %s\n", txid.String())
	users, err := db.GetUsers()
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func UpdateUser(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("UpdateUser | %s\n", txid.String())
	var update_user types.UserDto
	err := c.BodyParser(&update_user)
	if err != nil {
		log.Printf("Failed to parse user data\n%s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Failed to parse user data: %s\n", txid.String()))
	}
	current_user := c.Locals("user").(types.UserDbo)
	id, err := db.UpdateUser(current_user, update_user)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(id)
}
