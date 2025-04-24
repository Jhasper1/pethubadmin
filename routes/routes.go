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
	pethubRoutes.Get("/admin/getallpendingrequest", controllers.GetAllPendingRequests)
	pethubRoutes.Get("/admin/getalladopters", controllers.GetAllAdopters)
	pethubRoutes.Get("/admin/getallshelters", controllers.GetAllShelters)
	pethubRoutes.Post("/admin/updateregstatus", controllers.UpdateRegistrationStatus)
	pethubRoutes.Post("/admin/updateshelterstatus", controllers.UpdateShelterStatus)
	pethubRoutes.Post("/admin/updateadopterstatus", controllers.UpdateAdopterStatus)

	//try
	pethubRoutes.Get("/admin/getallshelterstry", controllers.GetAllSheltersAdmintry) // Route to get all shelters by id
	pethubRoutes.Put("/admin/shelters/:id/approve", controllers.ApproveShelterRegStatus)

	// ---------------- General Shared Routes ----------------
	app.Get("/allshelter", controllers.GetShelter)
	app.Get("/users/shelters/:id", controllers.GetAllSheltersByID)
}
