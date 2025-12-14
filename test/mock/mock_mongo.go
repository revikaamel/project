package mock

import (
	"context"
	"uas-backend/internal/model"

	"go.mongodb.org/mongo-driver/bson"
)

type MockMongoRepo struct {
	FindByIDFn        func(ctx context.Context, id string) (*model.AchievementMongo, error)
	CreateFn          func(ctx context.Context, ach *model.AchievementMongo) (string, error)
	UpdateFn          func(ctx context.Context, id string, data bson.M) error
	AddAttachmentsFn  func(ctx context.Context, id string, files []model.Attachment) error
	DeleteFn          func(ctx context.Context, id string) error
}

func (m *MockMongoRepo) FindByID(ctx context.Context, id string) (*model.AchievementMongo, error) {
	return m.FindByIDFn(ctx, id)
}
func (m *MockMongoRepo) Create(ctx context.Context, ach *model.AchievementMongo) (string, error) {
	return m.CreateFn(ctx, ach)
}
func (m *MockMongoRepo) Update(ctx context.Context, id string, data bson.M) error {
	return m.UpdateFn(ctx, id, data)
}
func (m *MockMongoRepo) AddAttachments(ctx context.Context, id string, files []model.Attachment) error {
	return m.AddAttachmentsFn(ctx, id, files)
}
func (m *MockMongoRepo) Delete(ctx context.Context, id string) error {
	return m.DeleteFn(ctx, id)
}
