package model

import "time"

type Student struct {
	ID        string `json:"id"`
	NIM       string `json:"nim"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	LecturerID string `json:"lecturer_id"` // dosen wali
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
