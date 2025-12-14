package model

import "time"

type AchievementRef struct {
    ID          string    `json:"id"`
    StudentID   string    `json:"student_id"`
    MongoID     string    `json:"mongo_id"`
    Status      string    `json:"status"` // draft, submitted, verified, rejected
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    // START PERUBAHAN KRUSIAL UNTUK MENGATASI ERROR NULL
    VerifiedBy  *string    `json:"verified_by,omitempty"` // MENGGUNAKAN POINTER (*)
    VerifiedAt  *time.Time `json:"verified_at,omitempty"` // MENGGUNAKAN POINTER (*)
    RejectedBy  *string    `json:"rejected_by,omitempty"` // MENGGUNAKAN POINTER (*)
    RejectedAt  *time.Time `json:"rejected_at,omitempty"` // MENGGUNAKAN POINTER (*)
    // END PERUBAHAN KRUSIAL
}