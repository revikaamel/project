package mock

import (
	"context"
	"uas-backend/internal/model"
)

type MockPGRepo struct {
	GetAllFn              func(ctx context.Context) ([]model.AchievementRef, error)
	GetByIDFn             func(ctx context.Context, id string) (*model.AchievementRef, error)
	CreateReferenceFn     func(ctx context.Context, studentID, mongoID string) (string, error)
	SetStatusSubmittedFn  func(ctx context.Context, id string) error
	SetStatusVerifiedFn   func(ctx context.Context, id, verifierID string) error
	SetStatusRejectedFn   func(ctx context.Context, id, verifierID string) error
	IsAdviseeOwnerFn      func(ctx context.Context, refID, lecturerID string) (bool, error)
}

func (m *MockPGRepo) GetAll(ctx context.Context) ([]model.AchievementRef, error) {
	return m.GetAllFn(ctx)
}
func (m *MockPGRepo) GetByID(ctx context.Context, id string) (*model.AchievementRef, error) {
	return m.GetByIDFn(ctx, id)
}
func (m *MockPGRepo) CreateReference(ctx context.Context, studentID, mongoID string) (string, error) {
	return m.CreateReferenceFn(ctx, studentID, mongoID)
}
func (m *MockPGRepo) SetStatusSubmitted(ctx context.Context, id string) error {
	return m.SetStatusSubmittedFn(ctx, id)
}
func (m *MockPGRepo) SetStatusVerified(ctx context.Context, id, verifierID string) error {
	return m.SetStatusVerifiedFn(ctx, refID, verifierID)
}
func (m *MockPGRepo) SetStatusRejected(ctx context.Context, id, verifierID string) error {
	return m.SetStatusRejectedFn(ctx, id, verifierID)
}
func (m *MockPGRepo) IsAdviseeOwner(ctx context.Context, refID, lecturerID string) (bool, error) {
	return m.IsAdviseeOwnerFn(ctx, refID, lecturerID)
}
