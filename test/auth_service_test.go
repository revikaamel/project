package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"uas-backend/internal/model"
	"uas-backend/internal/service"
	"uas-backend/config"
	testmock "uas-backend/test/mock"

	"github.com/gofiber/fiber/v2"
)

func TestLogin(t *testing.T) {
	mockUserRepo := &testmock.MockUserRepo{
		FindByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
			return &model.User{
				ID:       "user123",
				Email:    email,
				Password: "$2a$10$7uffrK/eEn6LfTS1sVvzBeuqHPwZPPjQH0zvSSiYjZafpqvbjR0Fe", // "password"
				Role:     "mahasiswa",
			}, nil
		},
	}

	cfg := &config.Config{
		JWTSecret:      "secret123",
		JWTExpireHours: 1,
	}

	svc := service.NewAuthService(mockUserRepo, cfg)

	app := fiber.New()
	app.Post("/login", svc.Login)

	body, _ := json.Marshal(map[string]string{
		"email":    "test@mail.com",
		"password": "password",
	})

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req)

	if res.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
}
