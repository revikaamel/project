package model

import "time"

type Lecturer struct {
	ID    string `json:"id"`
	NIDN  string `json:"nidn"`
	Name  string `json:"name"`
	Email string `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
