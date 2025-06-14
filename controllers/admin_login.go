package controllers

import (
	"errors"
	"pethubadmin/middleware"
	"pethubadmin/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginAdopter authenticates an adopter and retrieves their info

// ==============================================================

func LoginAdmin(c *fiber.Ctx) error {
	// Parse request body
	requestBody := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Check if the admin exists
	var adminAccount models.AdminAccount
	result := middleware.DBConn.Where("username = ?", requestBody.Username).First(&adminAccount)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	} else if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
		})
	}

	// Check password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(adminAccount.Password), []byte(requestBody.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(adminAccount.AdminID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generating token",
			"error":   err.Error(),
		})
	}

	// Login successful, return admin account and token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data": fiber.Map{
			"token":    token,
			"admin_id": adminAccount.AdminID, // Include admin ID in the response
			"admin":    adminAccount,
		},
	})
}
