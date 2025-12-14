package test

import (
	"context"
	"testing"

	"uas-backend/internal/model"
	"uas-backend/internal/service"
	testmock "uas-backend/test/mock"
)

func TestGetStudentByID(t *testing.T) {
	mockRepo := &testmock.MockStudentRepo{
		FindByIDFn: func(ctx context.Context, id string) (*model.Student, error) {
			return &model.Student{
				ID:   "S01",
				NIM:  "12345",
				Name: "Revika",
			}, nil
		},
	}

	svc := service.NewStudentService(mockRepo)

	student, err := svc.Repo.FindByID(context.Background(), "S01")
	if err != nil {
		t.Fatal(err)
	}

	if student.Name != "Revika" {
		t.Fatalf("expected Revika got %s", student.Name)
	}
}
