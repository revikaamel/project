package pg

import (
	"context"

	"uas-backend/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AchievementRefRepo struct {
	DB *pgxpool.Pool
}

func NewAchievementRefRepo(db *pgxpool.Pool) *AchievementRefRepo {
	return &AchievementRefRepo{DB: db}
}

// GetAll - Admin
func (r *AchievementRefRepo) GetAll(ctx context.Context) ([]model.AchievementRef, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, student_id, mongo_id, status, created_at, updated_at,
		 verified_by, verified_at, rejected_by, rejected_at
		 FROM achievement_refs ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.AchievementRef{}
	for rows.Next() {
		var a model.AchievementRef
		err := rows.Scan(
			&a.ID, &a.StudentID, &a.MongoID, &a.Status,
			&a.CreatedAt, &a.UpdatedAt,
			&a.VerifiedBy, &a.VerifiedAt, &a.RejectedBy, &a.RejectedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	return result, nil
}

// GetByStudentID - Mahasiswa
func (r *AchievementRefRepo) GetByStudentID(ctx context.Context, studentID string) ([]model.AchievementRef, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, student_id, mongo_id, status, created_at, updated_at,
		 verified_by, verified_at, rejected_by, rejected_at
		 FROM achievement_refs WHERE student_id=$1 ORDER BY created_at DESC`,
		studentID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.AchievementRef{}
	for rows.Next() {
		var a model.AchievementRef
		err := rows.Scan(
			&a.ID, &a.StudentID, &a.MongoID, &a.Status,
			&a.CreatedAt, &a.UpdatedAt,
			&a.VerifiedBy, &a.VerifiedAt, &a.RejectedBy, &a.RejectedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	return result, nil
}

// GetAdviseeAchievements - Dosen wali
func (r *AchievementRefRepo) GetAdviseeAchievements(ctx context.Context, lecturerID string) ([]model.AchievementRef, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT ar.id, ar.student_id, ar.mongo_id, ar.status, ar.created_at, ar.updated_at,
		        ar.verified_by, ar.verified_at, ar.rejected_by, ar.rejected_at
		   FROM achievement_refs ar
		   JOIN students s ON ar.student_id = s.id
		  WHERE s.lecturer_id = $1
		  ORDER BY ar.created_at DESC`,
		lecturerID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.AchievementRef{}
	for rows.Next() {
		var a model.AchievementRef
		err := rows.Scan(
			&a.ID, &a.StudentID, &a.MongoID, &a.Status,
			&a.CreatedAt, &a.UpdatedAt,
			&a.VerifiedBy, &a.VerifiedAt, &a.RejectedBy, &a.RejectedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	return result, nil
}

func (r *AchievementRefRepo) GetByID(ctx context.Context, id string) (*model.AchievementRef, error) {
	var a model.AchievementRef
	err := r.DB.QueryRow(ctx,
		`SELECT id, student_id, mongo_id, status, created_at, updated_at,
		 verified_by, verified_at, rejected_by, rejected_at
		 FROM achievement_refs WHERE id=$1`,
		id,
	).Scan(
		&a.ID, &a.StudentID, &a.MongoID, &a.Status,
		&a.CreatedAt, &a.UpdatedAt,
		&a.VerifiedBy, &a.VerifiedAt, &a.RejectedBy, &a.RejectedAt,
	)

	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AchievementRefRepo) CreateReference(ctx context.Context, studentID, mongoID string) (string, error) {
	var id string
	err := r.DB.QueryRow(ctx,
		`INSERT INTO achievement_refs
		  (student_id, mongo_id, status, created_at, updated_at)
		  VALUES ($1, $2, 'draft', NOW(), NOW())
		  RETURNING id`,
		studentID, mongoID,
	).Scan(&id)

	return id, err
}

func (r *AchievementRefRepo) SetStatusSubmitted(ctx context.Context, id string) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE achievement_refs
		    SET status='submitted',
		        updated_at=NOW()
		  WHERE id=$1`,
		id,
	)
	return err
}

func (r *AchievementRefRepo) SetStatusVerified(ctx context.Context, id, verifierID string) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE achievement_refs
		    SET status='verified',
		        verified_by=$1,
		        verified_at=NOW(),
		        updated_at=NOW()
		  WHERE id=$2`,
		verifierID, id,
	)
	return err
}

func (r *AchievementRefRepo) SetStatusRejected(ctx context.Context, id, verifierID string) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE achievement_refs
		    SET status='rejected',
		        rejected_by=$1,
		        rejected_at=NOW(),
		        updated_at=NOW()
		  WHERE id=$2`,
		verifierID, id,
	)
	return err
}

func (r *AchievementRefRepo) IsAdviseeOwner(ctx context.Context, refID, lecturerID string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1
			  FROM achievement_refs ar
			  JOIN students s ON ar.student_id = s.id
			 WHERE ar.id=$1 AND s.lecturer_id=$2
		)`,
		refID, lecturerID,
	).Scan(&exists)

	return exists, err
}

//softdelete
func (r *AchievementRefRepo) SetStatusDeleted(ctx context.Context, id string) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE achievement_refs
		    SET status='deleted',
		        updated_at=NOW()
		  WHERE id=$1`,
		id,
	)
	return err
}


