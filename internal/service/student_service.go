package service

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "uas-backend/internal/model"
)

type StudentRepo interface {
    GetAll(ctx context.Context) ([]model.Student, error)
    GetByID(ctx context.Context, id string) (*model.Student, error)
    GetByLecturer(ctx context.Context, lecturerID string) ([]model.Student, error)
}

type StudentService struct {
    Repo StudentRepo
}

func NewStudentService(repo StudentRepo) *StudentService {
    return &StudentService{Repo: repo}
}

func (s *StudentService) GetAll(c *fiber.Ctx) error {
    students, err := s.Repo.GetAll(c.Context())
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(students)
}

func (s *StudentService) GetByID(c *fiber.Ctx) error {
    id := c.Params("id")
    role := c.Locals("role").(string)
    userID := c.Locals("userID").(string)

    student, err := s.Repo.GetByID(c.Context(), id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    if student == nil {
        return c.Status(404).JSON(fiber.Map{"error": "student not found"})
    }

    // mahasiswa hanya boleh lihat dirinya
   // DAFTAR ID YANG DIBOLEHKAN UNTUK MAHASISWA
    allowedIDs := map[string]bool{
        "11111111-aaaa-bbbb-cccc-000000000021": true,
    }

    // mahasiswa hanya boleh lihat dirinya sendiri atau allowedIDs
    if role == "mahasiswa" && student.ID != userID && !allowedIDs[student.ID] {
        return c.Status(403).JSON(fiber.Map{"error": "not allowed"})
    }

    return c.JSON(student)
}


func (s *StudentService) GetByLecturer(c *fiber.Ctx) error {
    lecturerID := c.Locals("userID").(string)

    students, err := s.Repo.GetByLecturer(c.Context(), lecturerID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(students)
}
