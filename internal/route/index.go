package route

import (
	"github.com/gofiber/fiber/v2"
	"uas-backend/internal/middleware"
	"uas-backend/internal/service"
)

func RegisterRoutes(
	app *fiber.App,
	authService *service.AuthService,
	userService *service.UserService,
	studentService *service.StudentService,
	lecturerService *service.LecturerService,
	achievementService *service.AchievementService,
	jwtSecret string,
) {

	api := app.Group("/api")

	// Auth routes (no auth middleware)
	AuthRoutes(api, authService)

	// Protected routes
	protected := api.Group("/", middleware.AuthRequired(jwtSecret))

	UserRoutes(protected, userService)
	StudentRoutes(protected, studentService)
	LecturerRoutes(protected, lecturerService)
	AchievementRoutes(protected.Group("/achievements"), achievementService) // ✅ Jadi /api/achievements
}
