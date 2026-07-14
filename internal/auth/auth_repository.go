package auth

import (
	"context"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthRepository struct {
	queries *db.Queries
}

func NewAuthRepository(q *db.Queries) *AuthRepository {
	return &AuthRepository{
		queries: q,
	}
}
func (r *AuthRepository) CreateUser(
	ctx context.Context,
	params db.CreateUserParams,
) error {

	return r.queries.CreateUser(ctx, params)
}
func (r *AuthRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*db.User, error) {

	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, id pgtype.UUID) (*db.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*db.User, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
