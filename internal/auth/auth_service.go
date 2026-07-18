package auth

import (
	"context"

	"regexp"

	"github.com/Mikhail-Tal63/Orbit/configs"
	"github.com/Mikhail-Tal63/Orbit/internal/auth/errors"
	"github.com/Mikhail-Tal63/Orbit/internal/auth/validator"
	"github.com/Mikhail-Tal63/Orbit/internal/db"
	"github.com/Mikhail-Tal63/Orbit/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var usernameRe = regexp.MustCompile(`^[a-z0-9_]{3,20}$`)

type AuthService struct {
	authRepository AuthRepository
}

func NewAuthService(authRepository AuthRepository) *AuthService {
	return &AuthService{
		authRepository: authRepository,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, user *RegisterRequest) (*AuthResponce, error) {
	username := validator.NormalizeUsername(user.Username)
	email := validator.NormalizeEmail(user.Email)
	name := validator.NormalizeName(user.FirstName)

	if name == "" {
		name = username
	}

	if err := validator.ValidateUsername(username); err != nil {
		return nil, err
	}

	if err := validator.ValidateEmail(email); err != nil {
		return nil, err
	}

	if err := validator.ValidatePassword(user.Password); err != nil {
		return nil, err
	}

	if err := validator.ValidateName(name); err != nil {
		return nil, err
	}

	existingEmail, err := s.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, errors.ErrEmailAlreadyExists
	}

	existingUsername, err := s.authRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existingUsername != nil {
		return nil, errors.ErrUsernameTaken
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	params := db.CreateUserParams{
		ID:           uuid.New(),
		FirstName:    name,
		LastName:     validator.NormalizeName(user.LastName),
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		ImageID:      pgtype.UUID{Valid: false},
	}

	createdUser, err := s.authRepository.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	var imageID *uuid.UUID
	if createdUser.ImageID.Valid {
		id := uuid.UUID(createdUser.ImageID.Bytes)
		imageID = &id
	}

	secret := []byte(configs.Load().JWTSecret)

	accessToken, err := utils.CreateJWT(secret, createdUser.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(secret, createdUser.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponce{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserDTO{
			ID:        createdUser.ID,
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
			Username:  createdUser.Username,
			Email:     createdUser.Email,
			Phone:     createdUser.Phone,
			Role:      createdUser.Role,
			ImageID:   imageID,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, email, hashedpasswor string) (*AuthResponce, error) {
	user, err := s.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}
	if !utils.ComparePassword(user.PasswordHash, []byte(hashedpasswor)) {
		return nil, errors.ErrInvalidCredentials
	}

	secret := []byte(configs.Load().JWTSecret)
	token, err := utils.CreateJWT(secret, user.ID)
	if err != nil {
		return nil, err
	}

	refreshTokend, err := utils.GenerateRefreshToken(secret, user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponce{
		RefreshToken: refreshTokend,
		AccessToken:  token,
		User:         mapUserToDTO(user),
	}, nil
}
func mapUserToDTO(u *db.User) UserDTO {
	var imageID *uuid.UUID
	if u != nil && u.ImageID.Valid {
		id := uuid.UUID(u.ImageID.Bytes)
		imageID = &id
	}

	return UserDTO{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		ImageID:   imageID,
	}
}
