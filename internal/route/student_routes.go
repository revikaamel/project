package route

import (
	"github.com/gofiber/fiber/v2"
	"uas-backend/internal/middleware"
	"uas-backend/internal/service"
)

func StudentRoutes(router fiber.Router, studentService *service.StudentService) {
	std := router.Group("/students")

	std.Get("/", middleware.AdminOnly(), studentService.GetAll)
	std.Get("/:id", studentService.GetByID)
}
