package service

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "uas-backend/internal/model"
)

type LecturerRepo interface {
    GetAll(ctx context.Context) ([]model.Lecturer, error)
    GetByID(ctx context.Context, id string) (*model.Lecturer, error)
}

type LecturerService struct {
    Repo LecturerRepo
}

func NewLecturerService(repo LecturerRepo) *LecturerService {
    return &LecturerService{Repo: repo}
}

func (s *LecturerService) GetAll(c *fiber.Ctx) error {
    data, err := s.Repo.GetAll(c.Context())
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(data)
}

func (s *LecturerService) GetByID(c *fiber.Ctx) error {
    id := c.Params("id")

    data, err := s.Repo.GetByID(c.Context(), id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    if data == nil {
        return c.Status(404).JSON(fiber.Map{"error": "lecturer not found"})
    }

    return c.JSON(data)
}
