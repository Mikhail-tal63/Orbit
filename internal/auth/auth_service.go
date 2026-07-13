package auth

import (
	"context"

)

type AuthService struct{
	authRepository  AuthRepository
}

func NewAuthService(authRepository AuthRepository)*AuthService{
return &AuthService{
	authRepository: authRepository,
}
}

func (s *AuthService) CreateUser(ctx context.Context,payload *RegisterRequest) (error){
	return nil


}