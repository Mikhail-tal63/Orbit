package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)


// ==========================
// Mock Repository
// ==========================

type MockAuthRepository struct {
	mock.Mock
}


func (m *MockAuthRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*db.User, error) {

	args := m.Called(ctx, email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*db.User), args.Error(1)
}


func (m *MockAuthRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (*db.User, error) {

	args := m.Called(ctx, username)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*db.User), args.Error(1)
}

func (m *MockAuthRepository) CreateUser(
	ctx context.Context,
	params db.CreateUserParams,
) (*db.User, error) {

	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*db.User), args.Error(1)
}


func (m *MockAuthRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (*db.User, error) {

	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*db.User), args.Error(1)
}


func validRegisterRequest() *RegisterRequest {

	return &RegisterRequest{
		FirstName: "Jehad",
		LastName:  "Mohamed",
		Username:  "jehad",
		Email:     "jehad@test.com",
		Password:  "password123",
	}
}


func fakeUser() db.User {

	return db.User{
		ID: uuid.New(),

		FirstName: "Jehad",
		LastName:  "Mohamed",

		Username: "jehad",

		Email: "jehad@test.com",

		Phone: "",

		Role: "passenger",

		ImageID: pgtype.UUID{
			Valid: false,
		},
	}
}


func TestCreateUser_Success(t *testing.T) {


	repo := new(MockAuthRepository)


	user := fakeUser()


	repo.
		On(
			"GetUserByEmail",
			mock.Anything,
			"jehad@test.com",
		).
		Return(nil, nil)


	repo.
		On(
			"GetUserByUsername",
			mock.Anything,
			"jehad",
		).
		Return(nil, nil)


	repo.
		On(
			"CreateUser",
			mock.Anything,
			mock.Anything,
		).
		Return(&user, nil)



	service := NewAuthService(repo)



	res, err := service.CreateUser(
		context.Background(),
		validRegisterRequest(),
	)



	require.NoError(t, err)

	require.NotNil(t, res)


	require.NotEmpty(
		t,
		res.AccessToken,
	)


	require.NotEmpty(
		t,
		res.RefreshToken,
	)


	require.Equal(
		t,
		"jehad",
		res.User.Username,
	)


	require.Equal(
		t,
		"jehad@test.com",
		res.User.Email,
	)


	repo.AssertExpectations(t)

}




func TestCreateUser_EmailAlreadyExists(t *testing.T) {


	repo := new(MockAuthRepository)


	repo.
		On(
			"GetUserByEmail",
			mock.Anything,
			"jehad@test.com",
		).
		Return(&db.User{}, nil)



	service := NewAuthService(repo)



	res, err := service.CreateUser(
		context.Background(),
		validRegisterRequest(),
	)



	require.Error(t, err)

	require.Nil(t, res)


	repo.AssertExpectations(t)

}




func TestCreateUser_UsernameAlreadyExists(t *testing.T) {


	repo := new(MockAuthRepository)


	repo.
		On(
			"GetUserByEmail",
			mock.Anything,
			"jehad@test.com",
		).
		Return(nil, nil)


	repo.
		On(
			"GetUserByUsername",
			mock.Anything,
			"jehad",
		).
		Return(&db.User{}, nil)



	service := NewAuthService(repo)



	res, err := service.CreateUser(
		context.Background(),
		validRegisterRequest(),
	)



	require.Error(t, err)

	require.Nil(t, res)


	repo.AssertExpectations(t)

}




func TestCreateUser_InvalidInput(t *testing.T) {


	tests := []struct {
		name string
		user *RegisterRequest
	}{
		{
			name: "invalid username",
			user: &RegisterRequest{
				Username: "!!",
				Email: "test@test.com",
				Password: "password123",
			},
		},

		{
			name: "invalid email",
			user: &RegisterRequest{
				Username: "jehad",
				Email: "wrong",
				Password: "password123",
			},
		},

		{
			name: "weak password",
			user: &RegisterRequest{
				Username: "jehad",
				Email: "test@test.com",
				Password: "1",
			},
		},
	}


	for _, tt := range tests {


		t.Run(
			tt.name,
			func(t *testing.T){

				repo := new(MockAuthRepository)


				service := NewAuthService(repo)


				res, err := service.CreateUser(
					context.Background(),
					tt.user,
				)


				require.Error(t, err)

				require.Nil(t, res)

			},
		)
	}
}




func TestCreateUser_RepositoryError(t *testing.T) {


	repo := new(MockAuthRepository)



	dbError := errors.New("database error")



	repo.
		On(
			"GetUserByEmail",
			mock.Anything,
			"jehad@test.com",
		).
		Return(nil, dbError)



	service := NewAuthService(repo)



	res, err := service.CreateUser(
		context.Background(),
		validRegisterRequest(),
	)



	require.Error(t, err)

	require.Nil(t, res)


	require.ErrorIs(
		t,
		err,
		dbError,
	)


	repo.AssertExpectations(t)

}