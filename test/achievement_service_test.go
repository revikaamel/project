package test

import (
	"context"
	"net/http/httptest"
	"testing"

	"uas-backend/internal/model"
	"uas-backend/internal/service"
	testmock "uas-backend/test/mock"

	"github.com/gofiber/fiber/v2"
)

func TestSubmitAchievement(t *testing.T) {
	mockPG := &testmock.MockPGRepo{
		GetByIDFn: func(ctx context.Context, id string) (*model.AchievementRef, error) {
			return &model.AchievementRef{
				ID:        id,
				StudentID: "student123",
				Status:    "draft",
			}, nil
		},
		SetStatusSubmittedFn: func(ctx context.Context, id string) error {
			return nil
		},
	}

	mockMongo := &testmock.MockMongoRepo{}

	svc := service.NewAchievementService(mockMongo, mockPG, "./uploads")

	app := fiber.New()
	app.Post("/submit/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", "student123")
		c.Locals("role", "mahasiswa")
		return svc.Submit(c)
	})

	req := httptest.NewRequest("POST", "/submit/ach123", nil)
	res, _ := app.Test(req)

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
}
