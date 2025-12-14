package route

import (
	"github.com/gofiber/fiber/v2"
	"uas-backend/internal/middleware"
	"uas-backend/internal/service"
)

func UserRoutes(router fiber.Router, userService *service.UserService) {
	user := router.Group("/users")

	user.Get("/", middleware.AdminOnly(), userService.GetAll)
	user.Get("/:id", middleware.AdminOnly(), userService.GetByID)
	user.Post("/", middleware.AdminOnly(), userService.Create)
	user.Put("/:id", middleware.AdminOnly(), userService.Update)
	user.Delete("/:id", middleware.AdminOnly(), userService.Delete)
}
