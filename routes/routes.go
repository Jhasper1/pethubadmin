package routes

import (
	"pethubadmin/controllers"
	"pethubadmin/middleware"

	"github.com/gofiber/fiber/v2"
)

func AppRoutes(app *fiber.App) {

	pethubRoutes := app.Group("/api", middleware.JWTMiddleware())

	// ---------------- Admin Routes ----------------
	app.Post("/admin/register", controllers.RegisterAdmin)
	app.Post("/admin/login", controllers.LoginAdmin)
	app.Get("/admin/getallpendingrequest", controllers.GetAllPendingRequests)
	app.Get("/admin/getalladopters", controllers.GetAllAdopters)
	app.Get("/admin/getallshelters", controllers.GetAllShelters)
	app.Post("/admin/updateregstatus", controllers.UpdateRegistrationStatus)
	app.Post("/admin/updateshelterstatus", controllers.UpdateShelterStatus)
	app.Post("/admin/updateadopterstatus", controllers.UpdateAdopterStatus)

	//try
	pethubRoutes.Get("/admin/getallshelterstry", controllers.GetAllSheltersAdmintry) // Route to get all shelters by id
	app.Get("/admin/getalladopterstry", controllers.GetAllAdoptersAdmintry)          // Route to get all adopters by id
	pethubRoutes.Put("/admin/shelters/:id/approve", controllers.ApproveShelterRegStatus)
	app.Get("/admin/shelters/count", controllers.CountActiveShelters)
	app.Get("/admin/adopters/count", controllers.CountAdopters)
	app.Get("/admin/pets/count", controllers.CountPets)
	app.Get("/admin/adoptedpets/count", controllers.CountAdoptedPets)
	app.Get("/admin/pendingshelters/count", controllers.CountPendingShelters)
	app.Get("/admin/approvedshelters/count", controllers.CountApprovedShelters)
	pethubRoutes.Get("/admin/blockedadopters", controllers.GetInactiveAdopters)
	pethubRoutes.Put("/admin/adopters/:id/activate", controllers.ActivateAdopter)
	app.Get("/admin/allreports", controllers.GetSubmittedReports)
	app.Get("/admin/shelterpetcounts", controllers.GetShelterPetCounts)
	app.Get("/admin/vaccinecounts", controllers.GetShelterVaccinationCounts)
	app.Put("/admin/shelters/:id/status", controllers.UpdateShelterStatusByID)
	app.Get("/admin/blockedshelters", controllers.GetBlockedShelters)
	app.Put("/admin/shelters/:id/activate", controllers.UpdateShelterStatusByIDtoactive)
	app.Get("/applications/adopter/:adopter_id", controllers.GetApplicationsByAdopterID)
	//app.Get("/admin/adopter/:adopter_id", controllers.GetApplicationsByAdopterID)
	app.Get("/admin/notifications", controllers.GetAllNotifications)
	app.Get("/applications/shelter/:shelter_id", controllers.GetApplicationsByShelterID)
	app.Get("/shelter/:shelter_id/pets", controllers.GetPetsByShelterID)

	// ---------------- General Shared Routes ----------------
	pethubRoutes.Get("/allshelter", controllers.GetShelter)
	pethubRoutes.Get("/users/shelters/:id", controllers.GetAllSheltersByID)
}
