package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (
			id,
			first_name,
			last_name,
			email,
			password_hash,
			phone,
			role,
			image_id,
			is_active,
			created_at,
			updated_at,
			last_login_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12
		)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.Phone,
		user.Role,
		user.ImageID,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
		user.LastLoginAt,
	)

	if err != nil {
		return err
	}

	return nil
}
