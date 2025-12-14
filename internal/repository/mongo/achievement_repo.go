package mongo

import (
	"context"

	"uas-backend/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type AchievementMongoRepo struct {
	Collection *mongo.Collection
}

func NewAchievementMongoRepo(db *mongo.Database) *AchievementMongoRepo {
	return &AchievementMongoRepo{
		Collection: db.Collection("achievements"),
	}
}

// ░█▀▀█ █░░ ░█▀▀█ █▀▀█ █▀▀▄ █▀▀█ █▀▀
// CREATE
func (r *AchievementMongoRepo) Create(ctx context.Context, ach *model.AchievementMongo) (primitive.ObjectID, error) {
	ach.ID = primitive.NewObjectID()

	result, err := r.Collection.InsertOne(ctx, ach)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// ░█▀▀▀ █░░ █▀▀ █▀▀ ▀▀█▀▀ █▀▀
// READ (GET BY ObjectID)
func (r *AchievementMongoRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*model.AchievementMongo, error) {
	var ach model.AchievementMongo

	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&ach)
	if err != nil {
		return nil, err
	}

	return &ach, nil
}

// ░█▀▀▀ █░░ █▀▀ █▀▀ ▀▀█▀▀ █▀▀
// GET ALL by student (Mahasiswa)
func (r *AchievementMongoRepo) GetByStudent(ctx context.Context, studentID string) ([]model.AchievementMongo, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []model.AchievementMongo
	for cursor.Next(ctx) {
		var a model.AchievementMongo
		if err := cursor.Decode(&a); err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	return result, nil
}

// ░█▀▀▀ █▀▀█ █░░ █▀▀ █░░ █▀▀
// UPDATE MAIN DATA ONLY
func (r *AchievementMongoRepo) Update(ctx context.Context, id primitive.ObjectID, ach *model.AchievementMongo) error {
	update := bson.M{
		"$set": bson.M{
			"type":        ach.Type,
			"title":       ach.Title,
			"description": ach.Description,
			"details":     ach.Details,
			"attachments": ach.Attachments,
		},
	}

	_, err := r.Collection.UpdateByID(ctx, id, update)
	return err
}

// ░█░░█ ▀▀█▀▀ █▀▀█ █▀▀█ █▄░█ █▀▀
// DELETE
// SOFT DELETE (recommended by SRS)
func (r *AchievementMongoRepo) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"status":     "deleted",
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err := r.Collection.UpdateByID(ctx, id, update)
	return err
}
