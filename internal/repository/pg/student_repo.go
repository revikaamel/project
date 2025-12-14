package pg

import (
	"context"

	"uas-backend/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepo struct {
	DB *pgxpool.Pool
}

func NewStudentRepo(db *pgxpool.Pool) *StudentRepo {
	return &StudentRepo{DB: db}
}

func (r *StudentRepo) GetAll(ctx context.Context) ([]model.Student, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, nim, name, email, lecturer_id FROM students`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.Student{}
	for rows.Next() {
		var s model.Student
		err := rows.Scan(&s.ID, &s.NIM, &s.Name, &s.Email, &s.LecturerID)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (r *StudentRepo) GetByID(ctx context.Context, id string) (*model.Student, error) {
	var s model.Student

	err := r.DB.QueryRow(ctx,
		`SELECT id, nim, name, email, lecturer_id
		   FROM students
		  WHERE id=$1`,
		id,
	).Scan(&s.ID, &s.NIM, &s.Name, &s.Email, &s.LecturerID)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *StudentRepo) GetByLecturer(ctx context.Context, lecturerID string) ([]model.Student, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, nim, name, email, lecturer_id
		   FROM students
		  WHERE lecturer_id=$1`,
		lecturerID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.Student{}
	for rows.Next() {
		var s model.Student
		err := rows.Scan(&s.ID, &s.NIM, &s.Name, &s.Email, &s.LecturerID)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}
