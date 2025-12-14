package route

import (
	"github.com/gofiber/fiber/v2"
	"uas-backend/internal/middleware"
	"uas-backend/internal/service"
)

func LecturerRoutes(router fiber.Router, lecturerService *service.LecturerService) {
	lect := router.Group("/lecturers")

	lect.Get("/", middleware.AdminOnly(), lecturerService.GetAll)
	lect.Get("/:id", lecturerService.GetByID)
}
