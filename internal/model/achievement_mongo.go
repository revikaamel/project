package model

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementMongo struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    StudentID   string             `bson:"student_id" json:"student_id"`
    Type        string             `bson:"type" json:"type"`
    Title       string             `bson:"title" json:"title"`
    Description string             `bson:"description" json:"description"`
    Details     string             `bson:"details" json:"details"`
    Attachments []Attachment       `bson:"attachments" json:"attachments"`
    CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
    Status      string             `bson:"status" json:"status"` // <--- WAJIB TAMBAH
}
