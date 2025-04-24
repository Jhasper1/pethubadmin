package middleware

import (
	"fmt"
	"log"
	"pethubadmin/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
	DBErr  error
)

// ConnectDB initializes the connection to the PostgreSQL database using
// environment variables for configuration and assigns the connection
// to the global variable DBConn. It returns true if there was an error
// establishing the connection, otherwise false.
func ConnectDB() bool {
	// Database Config
	dns := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s TimeZone=%s",
		GetEnv("DB_HOST"), GetEnv("DB_PORT"), GetEnv("DB_NAME"),
		GetEnv("DB_USER"), GetEnv("DB_PASSWORD"), GetEnv("DB_SSLM"),
		GetEnv("DB_TMEZ"))

	DBConn, DBErr = gorm.Open(postgres.Open(dns), &gorm.Config{})
	if DBErr != nil {
		log.Printf("Database connection error: %v\n", DBErr)
		return true
	}

	log.Println("Database connection established successfully")

	// Auto-migrate models
	if err := DBConn.AutoMigrate(
		&models.AdopterAccount{},
		&models.AdopterInfo{},
		&models.ShelterAccount{},
		&models.ShelterInfo{},
		&models.ShelterMedia{},
	); err != nil {
		log.Printf("AutoMigrate error: %v\n", err)
		return true
	}

	return false
}

// Note: Remember to close the database connection when the application shuts down
// using DBConn.Close() if applicable.
