package route

import (
	"github.com/gofiber/fiber/v2"
	"uas-backend/internal/service"
)

func AuthRoutes(router fiber.Router, authService *service.AuthService) {
	auth := router.Group("/auth")

	auth.Post("/login", authService.Login)
}
