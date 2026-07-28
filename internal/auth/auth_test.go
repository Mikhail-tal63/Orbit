package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)


type FakeAuthRepository struct {
	usersByEmail    map[string]*db.User
	usersByUsername map[string]*db.User
	usersByID       map[uuid.UUID]*db.User
	errToReturn     error
}

func NewFakeAuthRepository() *FakeAuthRepository {
	return &FakeAuthRepository{
		usersByEmail:    make(map[string]*db.User),
		usersByUsername: make(map[string]*db.User),
		usersByID:       make(map[uuid.UUID]*db.User),
	}
}

func (f *FakeAuthRepository) CreateUser(ctx context.Context, params db.CreateUserParams) (*db.User, error) {
	if f.errToReturn != nil {
		return nil, f.errToReturn
	}

	user := &db.User{
		ID:           params.ID,
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
		ImageID:      params.ImageID,
		Role:         "passenger",
	}

	f.usersByEmail[params.Email] = user
	f.usersByUsername[params.Username] = user
	f.usersByID[params.ID] = user

	return user, nil
}

func (f *FakeAuthRepository) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	if f.errToReturn != nil {
		return nil, f.errToReturn
	}
	user, exists := f.usersByEmail[email]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (f *FakeAuthRepository) GetUserByUsername(ctx context.Context, username string) (*db.User, error) {
	if f.errToReturn != nil {
		return nil, f.errToReturn
	}
	user, exists := f.usersByUsername[username]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (f *FakeAuthRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*db.User, error) {
	if f.errToReturn != nil {
		return nil, f.errToReturn
	}
	user, exists := f.usersByID[id]
	if !exists {
		return nil, nil
	}
	return user, nil
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

// CreateUser Tests ***********************************************************************************

func TestCreateUser_Success(t *testing.T) {
	repo := NewFakeAuthRepository()
	service := NewAuthService(repo, []byte("test-secret"))

	res, err := service.CreateUser(context.Background(), validRegisterRequest())

	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res.AccessToken)
	require.NotEmpty(t, res.RefreshToken)
	require.Equal(t, "jehad", res.User.Username)
	require.Equal(t, "jehad@test.com", res.User.Email)
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	repo := NewFakeAuthRepository()
	service := NewAuthService(repo, []byte("test-secret"))


	_, err := service.CreateUser(context.Background(), validRegisterRequest())
	require.NoError(t, err)

	
	res, err := service.CreateUser(context.Background(), validRegisterRequest())

	require.ErrorIs(t, err, ErrEmailAlreadyExists)
	require.Nil(t, res)
}

func TestCreateUser_UsernameAlreadyExists(t *testing.T) {
	repo := NewFakeAuthRepository()
	service := NewAuthService(repo, []byte("test-secret"))


	_, err := service.CreateUser(context.Background(), validRegisterRequest())
	require.NoError(t, err)

	req := validRegisterRequest()
	req.Email = "another@test.com"
	res, err := service.CreateUser(context.Background(), req)

	require.ErrorIs(t, err, ErrUsernameTaken)
	require.Nil(t, res)
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
				Email:    "test@test.com",
				Password: "password123",
			},
		},
		{
			name: "invalid email",
			user: &RegisterRequest{
				Username: "jehad",
				Email:    "wrong",
				Password: "password123",
			},
		},
		{
			name: "weak password",
			user: &RegisterRequest{
				Username: "jehad",
				Email:    "test@test.com",
				Password: "1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewFakeAuthRepository()
			service := NewAuthService(repo, []byte("test-secret"))

			res, err := service.CreateUser(context.Background(), tt.user)

			require.Error(t, err)
			require.Nil(t, res)
		})
	}
}

func TestCreateUser_RepositoryError(t *testing.T) {
	dbError := errors.New("database error")

	repo := NewFakeAuthRepository()
	repo.errToReturn = dbError

	service := NewAuthService(repo, []byte("test-secret"))

	res, err := service.CreateUser(context.Background(), validRegisterRequest())

	require.ErrorIs(t, err, dbError)
	require.Nil(t, res)
}

// Login Tests *****************************************************************************************

func TestLogin_Success(t *testing.T) {
	repo := NewFakeAuthRepository()
	service := NewAuthService(repo, []byte("test-secret"))


	_, err := service.CreateUser(context.Background(), validRegisterRequest())
	require.NoError(t, err)

	
	res, err := service.Login(context.Background(), "jehad@test.com", "password123")

	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res.AccessToken)
	require.NotEmpty(t, res.RefreshToken)
	require.Equal(t, "jehad", res.User.Username)
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := NewFakeAuthRepository()
	service := NewAuthService(repo, []byte("test-secret"))

	res, err := service.Login(context.Background(), "nonexistent@test.com", "password123")

	require.ErrorIs(t, err, ErrUserNotFound)
	require.Nil(t, res)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	repo := NewFakeAuthRepository()
	service := NewAuthService(repo, []byte("test-secret"))


	_, err := service.CreateUser(context.Background(), validRegisterRequest())
	require.NoError(t, err)

	res, err := service.Login(context.Background(), "jehad@test.com", "wrongpassword")

	require.ErrorIs(t, err, ErrInvalidCredentials)
	require.Nil(t, res)
}
