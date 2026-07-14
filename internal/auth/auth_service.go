package auth

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Mikhail-Tal63/Orbit/configs"
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

	username := strings.ToLower(strings.TrimSpace(user.Username))
	if !usernameRe.MatchString(username) {
		return nil, fmt.Errorf("username must be 3-20 chars, lowercase letters, numbers or underscore")
	}

	name := strings.TrimSpace(user.FirstName)
	if name == "" {
		name = username
	}
	email := strings.ToLower(strings.TrimSpace(user.Email))
	existedEmail, err := s.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existedEmail != nil {
		return nil, fmt.Errorf("user with %s already exists", user.Email)
	}

	existingByUsername, err := s.authRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existingByUsername != nil {
		return nil, fmt.Errorf("username %q is already taken", username)
	}

	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	params := db.CreateUserParams{
		ID:           uuid.New(),
		FirstName:    name,
		LastName:     user.LastName,
		Username:     username,
		Email:        email,
		PasswordHash: hashed,
		ImageID:      pgtype.UUID{Valid: false},
	}

	createduser, err := s.authRepository.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	var imageID *uuid.UUID

	if createduser.ImageID.Valid {
		id := uuid.UUID(createduser.ImageID.Bytes)
		imageID = &id
	}
	secret := []byte(configs.Load().JWTSecret)
	token, err := utils.CreateJWT(secret, params.ID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateRefreshToken(secret, params.ID)
	if err != nil {
		return nil, err
	}
	return &AuthResponce{
		AccessToken:  token,
		RefreshToken: refreshToken,
		User: UserDTO{
			FirstName: createduser.FirstName,
			ID:        createduser.ID,
			LastName:  createduser.LastName,
			Username:  createduser.Username,
			Email:     createduser.Email,
			Phone:     createduser.Phone,
			Role:      createduser.Role,
			ImageID:   imageID,
		},
	}, nil

}
