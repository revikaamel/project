package pg

import (
	"context"

	"uas-backend/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LecturerRepo struct {
	DB *pgxpool.Pool
}

func NewLecturerRepo(db *pgxpool.Pool) *LecturerRepo {
	return &LecturerRepo{DB: db}
}

func (r *LecturerRepo) GetAll(ctx context.Context) ([]model.Lecturer, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, nidn, name, email FROM lecturers`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.Lecturer{}
	for rows.Next() {
		var l model.Lecturer
		err := rows.Scan(&l.ID, &l.NIDN, &l.Name, &l.Email)
		if err != nil {
			return nil, err
		}
		result = append(result, l)
	}
	return result, nil
}

func (r *LecturerRepo) GetByID(ctx context.Context, id string) (*model.Lecturer, error) {
	var l model.Lecturer

	err := r.DB.QueryRow(ctx,
		`SELECT id, nidn, name, email
		   FROM lecturers
		  WHERE id=$1`,
		id,
	).Scan(&l.ID, &l.NIDN, &l.Name, &l.Email)

	if err != nil {
		return nil, err
	}

	return &l, nil
}
