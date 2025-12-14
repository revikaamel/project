package route

import (
	"uas-backend/internal/middleware"
	"uas-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

func AchievementRoutes(r fiber.Router, achievementService *service.AchievementService) {
	r.Get("/", middleware.MahasiswaOnly(), achievementService.GetAll) // sebelumnya StudentOnly
	r.Get("/:id", middleware.MahasiswaOnly(), achievementService.GetDetail)
	r.Post("/", middleware.MahasiswaOnly(), achievementService.Create)
	r.Put("/:id", middleware.MahasiswaOnly(), achievementService.Update)
	r.Post("/:id/submit", middleware.MahasiswaOnly(), achievementService.Submit)

	r.Post("/:id/verify", middleware.AdminOrLecturer(), achievementService.Verify) // sebelumnya LecturerOrAdminOnly
	r.Post("/:id/reject", middleware.AdminOrLecturer(), achievementService.Reject)
	r.Post("/:id/upload", middleware.MahasiswaOnly(), achievementService.UploadAttachment)
	
	// 🟢 SOFT DELETE
	r.Delete("/:id", middleware.MahasiswaOnly(), achievementService.SoftDeleteHandler)

}
