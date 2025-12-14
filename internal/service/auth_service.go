package service

import (
	"context"
	
	"uas-backend/config"
	"uas-backend/internal/model"
	"uas-backend/internal/util"
	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	UserRepo interface {
		FindByEmail(ctx context.Context, email string) (*model.User, error)
	}
	Cfg *config.Config
}

func NewAuthService(
	userRepo interface {
		FindByEmail(ctx context.Context, email string) (*model.User, error)
	},
	cfg *config.Config,
) *AuthService {
	return &AuthService{UserRepo: userRepo, Cfg: cfg}
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	user, err := s.UserRepo.FindByEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	if !util.CheckPassword(req.Password, user.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// INGAT: pastikan field config kamu benar
	token, err := util.GenerateToken(
    user.ID,
    user.Role,
    user.Email,
    s.Cfg.JWTSecret,
    s.Cfg.JWTExpireHours, // langsung INT
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": token, "user": user})
}
