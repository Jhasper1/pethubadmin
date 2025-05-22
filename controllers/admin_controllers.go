package controllers

import (
	"errors"
	"fmt"
	"math"
	"pethubadmin/middleware"
	"pethubadmin/models"
	"time"

	"pethubadmin/models/response"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterAdmin(c *fiber.Ctx) error {
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

	// Check if username exists
	var existingAdmin models.AdminAccount
	result := middleware.DBConn.Where("username = ?", requestBody.Username).First(&existingAdmin)
	if result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Username already exists",
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}

	// Create admin account
	adminAccount := models.AdminAccount{
		Username: requestBody.Username,
		Password: string(hashedPassword), // Store hashed password
	}

	// Save admin account to the database
	if err := middleware.DBConn.Create(&adminAccount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to register admin",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Admin registered successfully",
		"data":    adminAccount,
	})
}

//BLOCK OR UNBLOCK ACCOUNTS

// UPDATE SHELTER STATUS
// ==============================================================
func UpdateShelterStatus(c *fiber.Ctx) error {
	// Parse request body
	var requestBody struct {
		ShelterID uint   `json:"shelter_id"` // Must be exported (capitalized)
		Status    string `json:"status"`     // Must be exported (capitalized)
	}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validate status (case-insensitive check)
	status := strings.ToLower(requestBody.Status)
	if status != "active" && status != "inactive" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Status must be either 'active' or 'inactive'",
		})
	}

	// Update shelter status directly in database
	result := middleware.DBConn.Model(&models.ShelterAccount{}).
		Where("shelter_id = ?", requestBody.ShelterID).
		Update("status", status) // Use the lowercased version

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update shelter status",
			"error":   result.Error.Error(),
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Shelter not found",
		})
	}

	// Fetch updated shelter to return in response
	var shelter models.ShelterAccount
	if err := middleware.DBConn.First(&shelter, requestBody.ShelterID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch updated shelter",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Shelter status updated successfully",
		"data": fiber.Map{
			"shelter_id": shelter.ShelterID,
			"username":   shelter.Username,
			"status":     shelter.Status,
		},
	})
}

// UPDATE ADOPTER STATUS
// ==============================================================
func UpdateAdopterStatus(c *fiber.Ctx) error {
	// Parse request body
	requestBody := struct {
		AdopterID uint   `json:"adopter_id"`
		Status    string `json:"status"`
	}{}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validate status
	if requestBody.Status != "active" && requestBody.Status != "inactive" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Status must be either 'active' or 'inactive'",
		})
	}

	// Find adopter
	var adopter models.AdopterAccount
	if err := middleware.DBConn.First(&adopter, requestBody.AdopterID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Adopter not found",
		})
	}

	// Update status
	adopter.Status = requestBody.Status

	// Save changes
	if err := middleware.DBConn.Save(&adopter).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update adopter status",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Adopter status updated successfully",
		"data": fiber.Map{
			"adopter_id": adopter.AdopterID,
			"username":   adopter.Username,
			"status":     adopter.Status,
		},
	})
}

// GetPendingShelterRequests retrieves all shelters with reg_status = 'pending'
// @Summary Get pending shelter requests
// @Description Get list of all shelter accounts with pending registration

func GetAllPendingRequests(c *fiber.Ctx) error {
	var pendingShelters []models.ShelterAccount

	// Get shelters with pending registration
	if err := middleware.DBConn.
		Where("reg_status = ?", "pending").
		Find(&pendingShelters).
		Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch pending shelters",
			"error":   err.Error(),
		})
	}

	// Get additional info for each shelter
	var results []fiber.Map
	for _, shelter := range pendingShelters {
		var info models.ShelterInfo
		if err := middleware.DBConn.
			Where("shelter_id = ?", shelter.ShelterID).
			First(&info).
			Error; err == nil {
			results = append(results, fiber.Map{
				"account": shelter,
				"info":    info,
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Pending shelter requests retrieved",
		"count":   len(results),
		"data":    results,
	})
}

//UPDATE REGISTRATION STATUS
// ==============================================================

func UpdateRegistrationStatus(c *fiber.Ctx) error {
	// Parse request
	var request struct {
		ShelterID uint   `json:"shelter_id"`
		RegStatus string `json:"reg_status"` // "approved" or "rejected"

	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validate status
	if request.RegStatus != "approved" && request.RegStatus != "rejected" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Reg status must be 'approved' or 'rejected'",
		})
	}

	// Find shelter
	var shelter models.ShelterAccount
	if err := middleware.DBConn.
		Where("shelter_id = ?", request.ShelterID).
		First(&shelter).
		Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Shelter not found",
		})
	}

	// Prepare updates
	updates := map[string]interface{}{
		"reg_status": request.RegStatus,
	}

	// If approved, also activate the account
	if request.RegStatus == "approved" {
		updates["status"] = "active"
	}

	// Add admin notes if provided
	//

	// Save changes
	if err := middleware.DBConn.
		Model(&shelter).
		Updates(updates).
		Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update registration status",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Shelter registration status updated",
		"data": fiber.Map{
			"shelter_id": shelter.ShelterID,
			"username":   shelter.Username,
			"reg_status": request.RegStatus,
		},
	})
}

// GetAllReports retrieves all submitted reports with filtering options

// UpdateReportStatus changes the status of a submitted report
func UpdateReportStatus(c *fiber.Ctx) error {
	reportID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid report ID",
		})
	}

	var request struct {
		Status    string `json:"status" validate:"required,oneof=pending reviewed resolved"`
		AdminNote string `json:"admin_note,omitempty"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Find the report
	var report models.SubmittedReport
	if err := middleware.DBConn.First(&report, reportID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Report not found",
		})
	}

	// Update report status
	updates := map[string]interface{}{
		"status": request.Status,
	}

	if request.AdminNote != "" {
		updates["admin_note"] = request.AdminNote
	}

	if err := middleware.DBConn.Model(&report).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update report status",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Report status updated successfully",
		"report": fiber.Map{
			"id":     report.ID, // Changed from report.ReportID to report.ID
			"status": report.Status,
		},
	})
}
func GetAllAdopters(c *fiber.Ctx) error {
	var accounts []models.AdopterAccount
	if err := middleware.DBConn.Find(&accounts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch adopter accounts",
		})
	}

	var infos []models.AdopterInfo
	if err := middleware.DBConn.Find(&infos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch adopter info",
		})
	}

	infoMap := make(map[uint]models.AdopterInfo)
	for _, info := range infos {
		infoMap[info.AdopterID] = info
	}

	var combined []fiber.Map
	for _, account := range accounts {
		if info, ok := infoMap[account.AdopterID]; ok {
			combined = append(combined, fiber.Map{
				"adopter": account,
				"info":    info,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Adopters retrieved successfully",
		"data":    combined,
	})
}
func GetAllShelters(c *fiber.Ctx) error {
	var accounts []models.ShelterAccount
	if err := middleware.DBConn.Find(&accounts).Error; err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Failed to fetch shelter accounts",
			Data:    nil,
		})
	}

	var infos []models.ShelterInfo
	if err := middleware.DBConn.Find(&infos).Error; err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Failed to fetch shelter info",
			Data:    nil,
		})
	}

	// Create a map of ShelterID to ShelterInfo for faster lookup
	infoMap := make(map[uint]models.ShelterInfo)
	for _, info := range infos {
		infoMap[info.ShelterID] = info
	}

	// Combine data
	shelters := []fiber.Map{}
	for _, account := range accounts {
		if info, ok := infoMap[account.ShelterID]; ok {
			shelters = append(shelters, fiber.Map{
				"shelter": account,
				"info":    info,
			})
		}
	}

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "All shelters retrieved successfully",
		Data:    shelters,
	})
}

//try

func GetAllSheltersAdmintry(c *fiber.Ctx) error {
	var accounts []models.ShelterAccount
	if err := middleware.DBConn.Find(&accounts).Error; err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Failed to fetch shelter accounts",
			Data:    nil,
		})
	}

	var infos []models.ShelterInfo
	if err := middleware.DBConn.Find(&infos).Error; err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Failed to fetch shelter info",
			Data:    nil,
		})
	}

	// Create a map of ShelterID to ShelterInfo for faster lookup
	infoMap := make(map[uint]models.ShelterInfo)
	for _, info := range infos {
		infoMap[info.ShelterID] = info
	}

	// Combine data
	var combined []fiber.Map
	for _, account := range accounts {
		if info, ok := infoMap[account.ShelterID]; ok {
			combined = append(combined, fiber.Map{
				"shelter":    account,
				"info":       info,
				"reg_status": account.RegStatus, // Include reg_status from ShelterAccount
			})
		}
	}

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "All shelters retrieved successfully",
		Data:    combined,
	})
}

func ApproveShelterRegStatus(c *fiber.Ctx) error {
	// Get shelter_id from the URL
	shelterID := c.Params("id")

	// Fetch the shelter to check its current status
	var shelter models.ShelterAccount
	if err := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Shelter not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
			"error":   err.Error(),
		})
	}

	// Check if the shelter is already approved
	if shelter.RegStatus == "approved" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Shelter is already approved",
		})
	}

	// Check if the shelter is pending
	if shelter.RegStatus != "pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Shelter registration is not pending",
		})
	}

	// Update the shelter's reg_status to "approved"
	result := middleware.DBConn.Model(&models.ShelterAccount{}).
		Where("shelter_id = ?", shelterID).
		Update("reg_status", "approved")

	// Check for errors during the update
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to approve shelter registration",
			"error":   result.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Shelter registration approved successfully",
	})
}

func GetShelter(c *fiber.Ctx) error {
	// Fetch all shelter info
	var shelters []models.ShelterInfo
	result := middleware.DBConn.Find(&shelters)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving shelters",
		})
	}

	return c.JSON(shelters)
}

func GetAllSheltersByID(c *fiber.Ctx) error {
	ShelterID := c.Params("id")

	// Fetch shelter info by ID
	var ShelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Table("shelterinfo").Where("shelter_id = ?", ShelterID).First(&ShelterInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "shelter info not found",
		})
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
		})
	}

	// Fetch shelter account associated with the shelter info
	var shelterData models.ShelterInfo
	accountResult := middleware.DBConn.Where("shelter_id = ?", ShelterID).First(&shelterData)

	if accountResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch shelter info",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter info retrieved successfully",
		"data": fiber.Map{
			"info": ShelterInfo,
		},
	})
}

func CountActiveShelters(c *fiber.Ctx) error {
	var count int64

	// Count shelters with status = "active"
	if err := middleware.DBConn.Model(&models.ShelterAccount{}).
		Where("status = ?", "active").
		Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to count active shelters",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Active shelters counted successfully",
		"count":   count,
	})
}

func CountAdopters(c *fiber.Ctx) error {
	var count int64

	// Count all adopters
	if err := middleware.DBConn.Model(&models.AdopterAccount{}).Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to count adopters",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Adopters counted successfully",
		"count":   count,
	})
}

func CountPets(c *fiber.Ctx) error {
	var count int64

	// Count all pets
	if err := middleware.DBConn.Model(&models.PetInfo{}).Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to count pets",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Pets counted successfully",
		"count":   count,
	})
}

func CountPendingShelters(c *fiber.Ctx) error {
	var count int64

	if err := middleware.DBConn.Model(&models.ShelterAccount{}).
		Where("reg_status = ?", "pending").
		Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to count pending shelters",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Pendingx shelters counted successfully",
		"count":   count,
	})
}

func CountApprovedShelters(c *fiber.Ctx) error {
	var count int64

	if err := middleware.DBConn.Model(&models.ShelterAccount{}).
		Where("reg_status = ?", "approved").
		Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to count approve shelters",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Pendingx shelters counted successfully",
		"count":   count,
	})
}

func CountAdoptedPets(c *fiber.Ctx) error {
	var count int64

	// Count pets where adoption_status is "adopted"
	if err := middleware.DBConn.Model(&models.PetInfo{}).
		Where("status = ?", "adopted").
		Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to count adopted pets",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Adopted pets counted successfully",
		"count":   count,
	})
}

func GetInactiveAdopters(c *fiber.Ctx) error {
	var inactiveAdopters []models.AdopterAccount

	// Fetch adopter accounts with status = "inactive"
	if err := middleware.DBConn.
		Where("status = ?", "inactive").
		Find(&inactiveAdopters).
		Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch inactive adopter accounts",
			"error":   err.Error(),
		})
	}

	// Fetch corresponding adopter info
	var adopterInfos []models.AdopterInfo
	if err := middleware.DBConn.
		Where("adopter_id IN ?", getAdopterIDs(inactiveAdopters)).
		Find(&adopterInfos).
		Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch adopter info",
			"error":   err.Error(),
		})
	}

	// Map adopter info by AdopterID for easy lookup
	infoMap := make(map[uint]models.AdopterInfo)
	for _, info := range adopterInfos {
		infoMap[info.AdopterID] = info
	}

	// Combine adopter accounts with their info
	var combined []fiber.Map
	for _, adopter := range inactiveAdopters {
		if info, ok := infoMap[adopter.AdopterID]; ok {
			combined = append(combined, fiber.Map{
				"adopter": adopter,
				"info":    info,
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Inactive adopters retrieved successfully",
		"data":    combined,
	})
}

// Helper function to extract AdopterIDs from AdopterAccount slice
func getAdopterIDs(adopters []models.AdopterAccount) []uint {
	var ids []uint
	for _, adopter := range adopters {
		ids = append(ids, adopter.AdopterID)
	}
	return ids
}
func ActivateAdopter(c *fiber.Ctx) error {
	// Get adopter_id from the URL
	adopterID := c.Params("id")

	// Fetch the adopter to check its current status
	var adopter models.AdopterAccount
	if err := middleware.DBConn.Where("adopter_id = ?", adopterID).First(&adopter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Adopter not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
			"error":   err.Error(),
		})
	}

	// Check if the adopter is already active
	if adopter.Status == "active" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Adopter is already active",
		})
	}

	// Check if the adopter is inactive
	if adopter.Status != "inactive" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Adopter status is not inactive",
		})
	}

	// Update the adopter's status to "active"
	result := middleware.DBConn.Model(&models.AdopterAccount{}).
		Where("adopter_id = ?", adopterID).
		Update("status", "active")

	// Check for errors during the update
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to activate adopter",
			"error":   result.Error.Error(),
		})
	}

	// Return success response
	return c.JSON(fiber.Map{
		"message": "Adopter activated successfully",
		"data": fiber.Map{
			"adopter_id": adopter.AdopterID,
			"username":   adopter.Username,
			"status":     "active",
		},
	})
}

func GetSubmittedReports(c *fiber.Ctx) error {
	// Response structures
	type ReportedBy struct {
		AdopterID    uint   `json:"adopter_id"`
		AdopterName  string `json:"adopter_name"`
		AdopterEmail string `json:"adopter_email"`
	}

	type ReportDetail struct {
		ID          uint       `json:"id"`
		Reason      string     `json:"reason"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		CreatedAt   string     `json:"created_at"` // Changed from time.Time to string
		ReportedBy  ReportedBy `json:"reported_by"`
	}

	type ShelterReportResponse struct {
		ShelterID      uint           `json:"shelter_id"`
		ShelterName    string         `json:"shelter_name"`
		ShelterEmail   string         `json:"shelter_email"`
		ShelterStatus  string         `json:"shelter_status"`
		ShelterProfile string         `json:"shelter_profile"`
		TotalReports   int            `json:"total_reports"`
		Reports        []ReportDetail `json:"reports"`
	}

	// Helper function to clean strings
	cleanString := func(s string) string {
		s = strings.ReplaceAll(s, "{", "")
		s = strings.ReplaceAll(s, "}", "")
		s = strings.ReplaceAll(s, "\"", "")
		return strings.TrimSpace(s)
	}

	// Helper function to format time
	formatTime := func(t time.Time) string {
		return t.Format("01-02-2006 03:04 PM") // MM-DD-YYYY hh:mm AM/PM format
	}

	// Fetch reports with preloaded relationships
	var submittedReports []models.SubmittedReport
	if err := middleware.DBConn.
		Where("status = ?", "reported").
		Preload("Shelter").
		Preload("Adopter").
		Find(&submittedReports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch reports",
			"error":   err.Error(),
		})
	}

	// Create a map to store shelter information
	shelterInfoMap := make(map[uint]models.ShelterInfo)
	shelterReportsMap := make(map[uint][]ReportDetail)

	for _, report := range submittedReports {
		// Skip if essential relationships are missing
		if report.Shelter.ShelterID == 0 || report.Adopter.AdopterID == 0 {
			continue
		}

		// Store shelter info if not already stored
		if _, exists := shelterInfoMap[report.ShelterID]; !exists {
			shelterInfoMap[report.ShelterID] = report.Shelter
		}

		// Verify shelter is active
		var shelterAccount models.ShelterAccount
		if err := middleware.DBConn.
			Where("shelter_id = ? AND status = ?", report.ShelterID, "active").
			First(&shelterAccount).Error; err != nil {
			continue
		}

		// Create report detail with cleaned strings and formatted time
		detail := ReportDetail{
			ID:          report.ID,
			Reason:      cleanString(report.Reason),
			Description: cleanString(report.Description),
			Status:      report.Status,
			CreatedAt:   formatTime(report.CreatedAt), // Format the time here
			ReportedBy: ReportedBy{
				AdopterID:    report.AdopterID,
				AdopterName:  fmt.Sprintf("%s %s", report.Adopter.FirstName, report.Adopter.LastName),
				AdopterEmail: report.Adopter.Email,
			},
		}

		shelterReportsMap[report.ShelterID] = append(shelterReportsMap[report.ShelterID], detail)
	}

	// Prepare final response
	var response []ShelterReportResponse
	for shelterID, reports := range shelterReportsMap {
		shelterInfo := shelterInfoMap[shelterID]

		// Get shelter profile image
		var shelterMedia models.ShelterMedia
		middleware.DBConn.
			Where("shelter_id = ?", shelterID).
			First(&shelterMedia)

		response = append(response, ShelterReportResponse{
			ShelterID:      shelterID,
			ShelterName:    shelterInfo.ShelterName,
			ShelterEmail:   shelterInfo.ShelterEmail,
			ShelterStatus:  "active",
			ShelterProfile: shelterMedia.ShelterProfile,
			TotalReports:   len(reports),
			Reports:        reports,
		})
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}

func GetShelterPetCounts(c *fiber.Ctx) error {
	type ShelterPetCount struct {
		ShelterID       uint    `json:"shelter_id"`
		ShelterName     string  `json:"shelter_name"`
		TotalPets       int     `json:"total_pets"`
		Cats            int     `json:"cats"`
		Dogs            int     `json:"dogs"`
		Vaccinated      int     `json:"vaccinated"`
		Unvaccinated    int     `json:"unvaccinated"`
		VaccinationRate float64 `json:"vaccination_rate"`
	}

	result := make([]ShelterPetCount, 0)

	// 1. Get active and approved shelters
	var shelters []models.ShelterAccount
	if err := middleware.DBConn.
		Where("status = ? AND reg_status = ?", "active", "approved").
		Find(&shelters).
		Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Shelter pet counts retrieved successfully",
			"data":    result,
		})
	}

	for _, shelter := range shelters {
		// 2. Get shelter info
		var shelterInfo models.ShelterInfo
		if err := middleware.DBConn.
			Where("shelter_id = ?", shelter.ShelterID).
			First(&shelterInfo).
			Error; err != nil {
			continue
		}

		// 3. Get non-archived pets with media - using correct table name
		var pets []models.PetInfo
		if err := middleware.DBConn.
			Table("petinfo"). // Explicit table name
			Preload("PetMedia").
			Where("shelter_id = ? AND status != ?", shelter.ShelterID, "archived").
			Find(&pets).
			Error; err != nil {
			continue
		}

		// Count pets by type and vaccination status
		counts := struct {
			Total        int
			Cats         int
			Dogs         int
			Vaccinated   int
			Unvaccinated int
		}{Total: len(pets)}

		for _, pet := range pets {
			// Count by pet type
			petType := strings.ToLower(pet.PetType)
			if petType == "cat" {
				counts.Cats++
			} else if petType == "dog" {
				counts.Dogs++
			}

			// Count vaccination status
			if pet.PetMedia.PetVaccine != "" && strings.TrimSpace(pet.PetMedia.PetVaccine) != "" {
				counts.Vaccinated++
			} else {
				counts.Unvaccinated++
			}
		}

		// Calculate vaccination rate
		var rate float64
		if counts.Total > 0 {
			rate = float64(counts.Vaccinated) / float64(counts.Total) * 100
		}

		result = append(result, ShelterPetCount{
			ShelterID:       shelter.ShelterID,
			ShelterName:     shelterInfo.ShelterName,
			TotalPets:       counts.Total,
			Cats:            counts.Cats,
			Dogs:            counts.Dogs,
			Vaccinated:      counts.Vaccinated,
			Unvaccinated:    counts.Unvaccinated,
			VaccinationRate: math.Round(rate*100) / 100, // Round to 2 decimal places
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter pet and vaccination counts retrieved successfully",
		"data":    result,
	})
}

func GetShelterVaccinationCounts(c *fiber.Ctx) error {
	type ShelterVaccinationCount struct {
		ShelterID       uint    `json:"shelter_id"`
		ShelterName     string  `json:"shelter_name"`
		Vaccinated      int     `json:"vaccinated"`
		Unvaccinated    int     `json:"unvaccinated"`
		VaccinationRate float64 `json:"vaccination_rate"`
	}

	result := make([]ShelterVaccinationCount, 0)

	// Get active and approved shelters
	var shelters []models.ShelterAccount
	if err := middleware.DBConn.
		Where("status = ? AND reg_status = ?", "active", "approved").
		Find(&shelters).
		Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Shelter vaccination counts retrieved successfully",
			"data":    result,
		})
	}

	for _, shelter := range shelters {
		// Get shelter info
		var shelterInfo models.ShelterInfo
		if err := middleware.DBConn.
			Where("shelter_id = ?", shelter.ShelterID).
			First(&shelterInfo).
			Error; err != nil {
			continue
		}

		// Get non-archived pets with media - using explicit table name
		var pets []models.PetInfo
		if err := middleware.DBConn.
			Table("petinfo"). // Explicit table name
			Preload("PetMedia").
			Where("shelter_id = ? AND status != ?", shelter.ShelterID, "archived").
			Find(&pets).
			Error; err != nil {
			continue
		}

		// Count vaccination status
		counts := struct {
			Vaccinated   int
			Unvaccinated int
		}{}

		for _, pet := range pets {
			// Safely check PetMedia (in case Preload failed)
			if pet.PetMedia.PetVaccine != "" && strings.TrimSpace(pet.PetMedia.PetVaccine) != "" {
				counts.Vaccinated++
			} else {
				counts.Unvaccinated++
			}
		}

		// Calculate vaccination rate
		total := counts.Vaccinated + counts.Unvaccinated
		var rate float64
		if total > 0 {
			rate = float64(counts.Vaccinated) / float64(total) * 100
		}

		result = append(result, ShelterVaccinationCount{
			ShelterID:       shelter.ShelterID,
			ShelterName:     shelterInfo.ShelterName,
			Vaccinated:      counts.Vaccinated,
			Unvaccinated:    counts.Unvaccinated,
			VaccinationRate: math.Round(rate*100) / 100, // Round to 2 decimal places
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter vaccination counts retrieved successfully",
		"data":    result,
	})
}

func UpdateShelterStatusByID(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	id, err := strconv.ParseUint(shelterID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid shelter ID",
		})
	}

	// Check if shelter account exists
	var shelterAccount models.ShelterAccount
	if err := middleware.DBConn.Where("shelter_id = ?", uint(id)).First(&shelterAccount).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Shelter account not found",
			"error":   err.Error(),
		})
	}

	// Check if shelter info exists
	var shelterInfo models.ShelterInfo
	if err := middleware.DBConn.First(&shelterInfo, uint(id)).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Shelter info not found",
			"error":   err.Error(),
		})
	}

	// Update status in ShelterAccount
	result := middleware.DBConn.Model(&models.ShelterAccount{}).
		Where("shelter_id = ?", uint(id)).
		Update("status", "inactive")

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update shelter status",
			"error":   result.Error.Error(),
		})
	}

	// Ensure ALL reports with shelter_id and status 'reported' are updated to 'blocked'
	reportUpdate := middleware.DBConn.Model(&models.SubmittedReport{}).
		Where("shelter_id = ? AND status = ?", uint(id), "reported").
		Updates(map[string]interface{}{"status": "blocked"})

	if reportUpdate.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update submitted reports status",
			"error":   reportUpdate.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":         "Shelter blocked successfully",
		"shelter_id":      shelterAccount.ShelterID,
		"shelter_name":    shelterInfo.ShelterName,
		"shelter_email":   shelterInfo.ShelterEmail,
		"shelter_status":  "inactive",
		"reports_updated": reportUpdate.RowsAffected,
	})
}

func GetBlockedShelters(c *fiber.Ctx) error {
	// Response structures
	type ReportedBy struct {
		AdopterID    uint   `json:"adopter_id"`
		AdopterName  string `json:"adopter_name"`
		AdopterEmail string `json:"adopter_email"`
	}

	type ReportDetail struct {
		ID          uint       `json:"id"`
		Reason      string     `json:"reason"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		CreatedAt   string     `json:"created_at"` // Changed from time.Time to string
		ReportedBy  ReportedBy `json:"reported_by"`
	}

	type ShelterReportResponse struct {
		ShelterID      uint           `json:"shelter_id"`
		ShelterName    string         `json:"shelter_name"`
		ShelterEmail   string         `json:"shelter_email"`
		ShelterStatus  string         `json:"shelter_status"`
		ShelterProfile string         `json:"shelter_profile"`
		TotalReports   int            `json:"total_reports"`
		Reports        []ReportDetail `json:"reports"`
	}

	// Helper function to clean strings
	cleanString := func(s string) string {
		s = strings.ReplaceAll(s, "{", "")
		s = strings.ReplaceAll(s, "}", "")
		s = strings.ReplaceAll(s, "\"", "")
		return strings.TrimSpace(s)
	}

	// Helper function to format time
	formatTime := func(t time.Time) string {
		return t.Format("01-02-2006 03:04 PM") // MM-DD-YYYY hh:mm AM/PM format
	}

	// Fetch reports with preloaded relationships
	var submittedReports []models.SubmittedReport
	if err := middleware.DBConn.
		Where("status = ?", "blocked").
		Preload("Shelter").
		Preload("Adopter").
		Find(&submittedReports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch reports",
			"error":   err.Error(),
		})
	}

	// Create a map to store shelter information
	shelterInfoMap := make(map[uint]models.ShelterInfo)
	shelterReportsMap := make(map[uint][]ReportDetail)

	for _, report := range submittedReports {
		// Skip if essential relationships are missing
		if report.Shelter.ShelterID == 0 || report.Adopter.AdopterID == 0 {
			continue
		}

		// Store shelter info if not already stored
		if _, exists := shelterInfoMap[report.ShelterID]; !exists {
			shelterInfoMap[report.ShelterID] = report.Shelter
		}

		// Verify shelter is active
		var shelterAccount models.ShelterAccount
		if err := middleware.DBConn.
			Where("shelter_id = ? AND status = ?", report.ShelterID, "inactive").
			First(&shelterAccount).Error; err != nil {
			continue
		}

		// Create report detail with cleaned strings and formatted time
		detail := ReportDetail{
			ID:          report.ID,
			Reason:      cleanString(report.Reason),
			Description: cleanString(report.Description),
			Status:      report.Status,
			CreatedAt:   formatTime(report.CreatedAt), // Format the time here
			ReportedBy: ReportedBy{
				AdopterID:    report.AdopterID,
				AdopterName:  fmt.Sprintf("%s %s", report.Adopter.FirstName, report.Adopter.LastName),
				AdopterEmail: report.Adopter.Email,
			},
		}

		shelterReportsMap[report.ShelterID] = append(shelterReportsMap[report.ShelterID], detail)
	}

	// Prepare final response
	var response []ShelterReportResponse
	for shelterID, reports := range shelterReportsMap {
		shelterInfo := shelterInfoMap[shelterID]

		// Get shelter profile image
		var shelterMedia models.ShelterMedia
		middleware.DBConn.
			Where("shelter_id = ?", shelterID).
			First(&shelterMedia)

		response = append(response, ShelterReportResponse{
			ShelterID:      shelterID,
			ShelterName:    shelterInfo.ShelterName,
			ShelterEmail:   shelterInfo.ShelterEmail,
			ShelterStatus:  "inactive",
			ShelterProfile: shelterMedia.ShelterProfile,
			TotalReports:   len(reports),
			Reports:        reports,
		})
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}

func UpdateShelterStatusByIDtoactive(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	id, err := strconv.ParseUint(shelterID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid shelter ID",
		})
	}

	// Check if shelter account exists
	var shelterAccount models.ShelterAccount
	if err := middleware.DBConn.Where("shelter_id = ?", uint(id)).First(&shelterAccount).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Shelter account not found",
			"error":   err.Error(),
		})
	}

	// Check if shelter info exists
	var shelterInfo models.ShelterInfo
	if err := middleware.DBConn.First(&shelterInfo, uint(id)).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Shelter info not found",
			"error":   err.Error(),
		})
	}

	// Update status in ShelterAccount
	result := middleware.DBConn.Model(&models.ShelterAccount{}).
		Where("shelter_id = ?", uint(id)).
		Update("status", "active")

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update shelter status",
			"error":   result.Error.Error(),
		})
	}

	// Ensure ALL reports with shelter_id and status 'reported' are updated to 'blocked'
	reportUpdate := middleware.DBConn.Model(&models.SubmittedReport{}).
		Where("shelter_id = ? AND status = ?", uint(id), "blocked").
		Updates(map[string]interface{}{"status": "resolved"})

	if reportUpdate.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update submitted reports status",
			"error":   reportUpdate.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":         "Shelter blocked successfully",
		"shelter_id":      shelterAccount.ShelterID,
		"shelter_name":    shelterInfo.ShelterName,
		"shelter_email":   shelterInfo.ShelterEmail,
		"shelter_status":  "active",
		"reports_updated": reportUpdate.RowsAffected,
	})
}
