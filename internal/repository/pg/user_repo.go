package pg

import (
	"context"

	"uas-backend/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	DB *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) GetAll(ctx context.Context) ([]model.User, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, email, role FROM users`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.User{}
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Email, &u.Role)
		if err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User

	err := r.DB.QueryRow(ctx,
		`SELECT id, email, role 
		   FROM users 
		  WHERE id=$1`,
		id,
	).Scan(&u.ID, &u.Email, &u.Role)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User

	err := r.DB.QueryRow(ctx,
		`SELECT id, email, password, role 
		   FROM users 
		  WHERE email=$1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Password, &u.Role)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) (string, error) {
	var id string

	err := r.DB.QueryRow(ctx,
		`INSERT INTO users (email, password, role) 
		          VALUES ($1, $2, $3) 
		       RETURNING id`,
		user.Email, user.Password, user.Role,
	).Scan(&id)

	return id, err
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE users 
		    SET email=$1, password=$2, role=$3 
		  WHERE id=$4`,
		user.Email, user.Password, user.Role, user.ID,
	)
	return err
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	_, err := r.DB.Exec(ctx,
		`DELETE FROM users WHERE id=$1`,
		id,
	)
	return err
}
