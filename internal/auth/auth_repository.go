package auth

import (
	"context"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
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
