package test

import (
	"context"
	"testing"

	"uas-backend/internal/model"
	"uas-backend/internal/service"
	testmock "uas-backend/test/mock"
)

func TestGetLecturerByID(t *testing.T) {
	mockRepo := &testmock.MockLecturerRepo{
		FindByIDFn: func(ctx context.Context, id string) (*model.Lecturer, error) {
			return &model.Lecturer{
				ID:   "L001",
				Name: "Dosen Pembimbing",
			}, nil
		},
	}

	svc := service.NewLecturerService(mockRepo)

	lect, err := svc.Repo.FindByID(context.Background(), "L001")
	if err != nil {
		t.Fatal(err)
	}

	if lect.Name != "Dosen Pembimbing" {
		t.Fatalf("expected Dosen Pembimbing got %s", lect.Name)
	}
}
