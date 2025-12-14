package test

import (
	"context"
	"testing"

	"uas-backend/internal/model"
	"uas-backend/internal/service"
	testmock "uas-backend/test/mock"
)

func TestGetUserByID(t *testing.T) {
	mockRepo := &testmock.MockUserRepo{
		FindByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return &model.User{
				ID:    id,
				Email: "example@mail.com",
				Role:  "admin",
			}, nil
		},
	}

	svc := service.NewUserService(mockRepo)

	user, err := svc.Repo.FindByID(context.Background(), "user123")
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != "user123" {
		t.Fatalf("expected user123 got %s", user.ID)
	}
}
