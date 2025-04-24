package middleware

import (
	"fmt"
	"log"
	"os"
	"pethubadmin/models"
	"pethubadmin/models/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load() // Load from .env file
	if err != nil {
		log.Println("Warning: No .env file found")
	}
}

// Secret key for signing tokens (should be stored in env variables)
var SecretKey = os.Getenv("SECRET_KEY")

// GenerateJWT generates a new JWT token
func GenerateJWT(ID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // Expires in 72 hours

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")

		if tokenString == "" {
			return c.JSON(response.ResponseModel{
				RetCode: "401",
				Message: "Unauthorized: No token provided",
				Data:    nil,
			})
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// var count int64
		// err := DBConn.Table("token_blacklists").Where("token = ?", tokenString).Count(&count).Error
		// if err != nil {
		// 	return c.JSON(response.ResponseModel{
		// 		RetCode: "500",
		// 		Message: "Error checking token blacklist",
		// 		Data:    err,
		// 	})
		// }

		// if count > 0 {
		// 	return c.JSON(response.ResponseModel{
		// 		RetCode: "401",
		// 		Message: "Unauthorized: Token is blacklisted",
		// 		Data:    nil,
		// 	})
		// }

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(SecretKey), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(response.ResponseModel{
				RetCode: "401",
				Message: "Unauthorized: Invalid token",
				Data:    nil,
			})
		}

		fmt.Println("Token received:", tokenString)

		c.Locals("id", token.Claims.(jwt.MapClaims)["id"])
		return c.Next()
	}
}

func GetAllShelterUserHomeScreen(c *fiber.Ctx) error {
	// Define a struct to include shelter info and media fields
	type ShelterWithMedia struct {
		models.ShelterInfo
		ShelterProfile *string `json:"shelter_profile"`
		ShelterCover   *string `json:"shelter_cover"`
	}

	// Create a slice to store all shelters with their media
	var shelters []ShelterWithMedia

	// Fetch all shelters and their media from the database
	result := DBConn.Table("shelterinfo").
		Select("shelterinfo.*, sheltermedia.shelter_profile, sheltermedia.shelter_cover").
		Joins("LEFT JOIN sheltermedia ON shelterinfo.shelter_id = sheltermedia.shelter_id").
		Find(&shelters)

	// Handle errors
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving shelter data",
			"error":   result.Error.Error(),
		})
	}

	// Return the list of shelters with their media
	return c.JSON(fiber.Map{
		"message": "Shelters retrieved successfully",
		"data":    shelters,
	})
}
