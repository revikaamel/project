package middleware

import (
	"strings" // <-- TAMBAHKAN INI
	"github.com/gofiber/fiber/v2"
)

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
        // PERBAIKAN: Ubah role ke huruf kecil untuk perbandingan yang case-insensitive
		if strings.ToLower(role) != "admin" {
			return c.Status(403).JSON(fiber.Map{"error": "admin only"})
		}
		return c.Next()
	}
}

func MahasiswaOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
        // PERBAIKAN KRUSIAL DI SINI: Ubah role ke huruf kecil untuk perbandingan yang case-insensitive
		if strings.ToLower(role) != "mahasiswa" { 
			return c.Status(403).JSON(fiber.Map{"error": "students only"})
		}
		return c.Next()
	}
}

func AdminOrLecturer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
        // PERBAIKAN: Ubah role ke huruf kecil untuk perbandingan yang case-insensitive
        normalizedRole := strings.ToLower(role) 
		if normalizedRole != "admin" && normalizedRole != "dosen" {
			return c.Status(403).JSON(fiber.Map{"error": "admin or lecturer only"})
		}
		return c.Next()
	}
}