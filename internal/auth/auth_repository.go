package auth

import (
	"context"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
	"github.com/google/uuid"
)

type AuthRepository interface {

	CreateUser(
		ctx context.Context,
		params db.CreateUserParams,
	) (*db.User, error)


	GetUserByEmail(
		ctx context.Context,
		email string,
	) (*db.User, error)


	GetUserByUsername(
		ctx context.Context,
		username string,
	) (*db.User, error)


	GetUserByID(
		ctx context.Context,
		id uuid.UUID,
	) (*db.User, error)
}


type AuthRepositoryImpl struct {
	queries *db.Queries
}

var _ AuthRepository = (*AuthRepositoryImpl)(nil)



func NewAuthRepository(q *db.Queries) *AuthRepositoryImpl {

	return &AuthRepositoryImpl{
		queries: q,
	}
}


func (r *AuthRepositoryImpl) CreateUser(
	ctx context.Context,
	params db.CreateUserParams,
) (*db.User, error) {

	user, err := r.queries.CreateUser(ctx, params)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepositoryImpl) GetUserByEmail(
	ctx context.Context,
	email string,
) (*db.User, error) {

	user, err := r.queries.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}



func (r *AuthRepositoryImpl) GetUserByUsername(
	ctx context.Context,
	username string,
) (*db.User, error) {

	user, err := r.queries.GetUserByUsername(ctx, username)

	if err != nil {
		return nil, err
	}

	return &user, nil
}



func (r *AuthRepositoryImpl) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (*db.User, error) {

	user, err := r.queries.GetUserByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}