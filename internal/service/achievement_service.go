package service

import (
	"context"
	"os"
	"path/filepath"

	"uas-backend/internal/model"
	pgRepo "uas-backend/internal/repository/pg"
	mongoRepo "uas-backend/internal/repository/mongo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AchievementService struct {
	PgRepo    *pgRepo.AchievementRefRepo
	MongoRepo *mongoRepo.AchievementMongoRepo
}

func NewAchievementService(pg *pgRepo.AchievementRefRepo, mg *mongoRepo.AchievementMongoRepo) *AchievementService {
	return &AchievementService{
		PgRepo:    pg,
		MongoRepo: mg,
	}
}

//
// ─── GET ALL ──────────────────────────────────────────────────
// Admin → semua
// Dosen → semua mahasiswa bimbingannya
// Mahasiswa → hanya miliknya
//
func (s *AchievementService) GetAll(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("userID").(string)

	var data []model.AchievementRef
	var err error

	ctx := context.Background()

	switch role {
	case "admin":
		data, err = s.PgRepo.GetAll(ctx)

	case "lecturer":
		data, err = s.PgRepo.GetAdviseeAchievements(ctx, userID)

	case "mahasiswa":
		data, err = s.PgRepo.GetByStudentID(ctx, userID)

	default:
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized role"})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(data)
}

//
// ─── GET DETAIL (PG + MONGO) ───────────────────────────────────
//
func (s *AchievementService) GetDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx := context.Background()

	ref, err := s.PgRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	mongoID, err := primitive.ObjectIDFromHex(ref.MongoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "invalid mongo id"})
	}

	ach, err := s.MongoRepo.GetByID(ctx, mongoID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "achievement detail not found"})
	}

	return c.JSON(fiber.Map{
		"reference": ref,
		"detail":    ach,
	})
}

//
// ─── CREATE Achievement (PG Ref + Mongo Detail) ─────────────────
//
func (s *AchievementService) Create(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	studentID := c.Locals("userID").(string)

	if role != "mahasiswa" {
		return c.Status(403).JSON(fiber.Map{"error": "Only student can create achievement"})
	}

	var input model.AchievementMongo
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	input.StudentID = studentID

	ctx := context.Background()

	// Insert Mongo
	mongoID, err := s.MongoRepo.Create(ctx, &input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Insert reference in PG
	refID, err := s.PgRepo.CreateReference(ctx, studentID, mongoID.Hex())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"reference_id": refID,
		"mongo_id":     mongoID.Hex(),
	})
}

//
// ─── UPDATE ─────────────────────────────────────────────────────
// Mahasiswa hanya boleh update prestasi miliknya
//
func (s *AchievementService) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	role := c.Locals("role").(string)
	studentID := c.Locals("userID").(string)

	ctx := context.Background()

	ref, err := s.PgRepo.GetByID(ctx, id)
	if err != nil || ref.StudentID != studentID || role != "mahasiswa" {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized update"})
	}

	var input model.AchievementMongo
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoID)

	err = s.MongoRepo.Update(ctx, mongoID, &input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Updated successfully",
	})
}

//
// ─── SUBMIT ─────────────────────────────────────────────────────
// Mahasiswa → ubah status dari draft → submitted
//
func (s *AchievementService) Submit(c *fiber.Ctx) error {
	id := c.Params("id")

	role := c.Locals("role").(string)
	userID := c.Locals("userID").(string)

	if role != "mahasiswa" {
		return c.Status(403).JSON(fiber.Map{"error": "Only student can submit"})
	}

	ctx := context.Background()

	ref, err := s.PgRepo.GetByID(ctx, id)
	if err != nil || ref.StudentID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
	}

	err = s.PgRepo.SetStatusSubmitted(ctx, id)
	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Submitted"})
}

//
// ─── VERIFY / REJECT (Dosen Wali + Admin) ───────────────────────
//
func (s *AchievementService) Verify(c *fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("userID").(string)

	if role != "dosen" && role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
	}

	ctx := context.Background()

	if role == "lecturer" {
		ok, _ := s.PgRepo.IsAdviseeOwner(ctx, id, userID)
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "Not your student"})
		}
	}

	err := s.PgRepo.SetStatusVerified(ctx, id, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Verified"})
}

func (s *AchievementService) Reject(c *fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("userID").(string)

	if role != "dosen" && role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
	}

	ctx := context.Background()

	if role == "mahasiswa" {
		ok, _ := s.PgRepo.IsAdviseeOwner(ctx, id, userID)
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "Not your student"})
		}
	}

	err := s.PgRepo.SetStatusRejected(ctx, id, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Rejected"})
}

//
// ─── UPLOAD FILE (LOCAL + MONGO METADATA) ───────────────────────
//
func (s *AchievementService) UploadAttachment(c *fiber.Ctx) error {
	id := c.Params("id")

	role := c.Locals("role").(string)
	studentID := c.Locals("userID").(string)

	if role != "mahasiswa" {
		return c.Status(403).JSON(fiber.Map{"error": "Students only"})
	}

	ctx := context.Background()

	ref, err := s.PgRepo.GetByID(ctx, id)
	if err != nil || ref.StudentID != studentID {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// upload file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "File required"})
	}

	savePath := filepath.Join("uploads", file.Filename)
	os.MkdirAll("uploads", 0755)
	err = c.SaveFile(file, savePath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// update mongo metadata
	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoID)
	ach, _ := s.MongoRepo.GetByID(ctx, mongoID)

	ach.Attachments = append(ach.Attachments, model.Attachment{
    FileName:   file.Filename,
    URL:        "/uploads/" + file.Filename,
    MimeType:   file.Header.Get("Content-Type"),
    Size:       file.Size,
    UploadedAt: time.Now(),
	})


	err = s.MongoRepo.Update(ctx, mongoID, ach)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "File uploaded",
		"file":    file.Filename,
	})
}
//softdelete
func (s *AchievementService) SoftDelete(ctx context.Context, refID string) error {
    // ambil data ref → dapatkan MongoID
    ref, err := s.PgRepo.GetByID(ctx, refID)
    if err != nil {
        return err
    }

    // convert mongo_id string → ObjectID
    mongoID, err := primitive.ObjectIDFromHex(ref.MongoID)
    if err != nil {
        return err
    }

    // soft delete MongoDB
    if err := s.MongoRepo.SoftDelete(ctx, mongoID); err != nil {
        return err
    }

    // soft delete PG
    return s.PgRepo.SetStatusDeleted(ctx, refID)
}
func (s *AchievementService) SoftDeleteHandler(c *fiber.Ctx) error {
    id := c.Params("id")

    // panggil business logic dari SoftDelete()
    if err := s.SoftDelete(c.Context(), id); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "message": "Achievement deleted successfully",
    })
}
