package middleware

import (
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

func Recover() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC: %v\n%s", r, string(debug.Stack()))
				_ = c.Status(500).JSON(fiber.Map{"error": "internal server error"})
			}
		}()

		return c.Next()
	}
}
