package auth

import (
	"context"
	"time"
)

type AuthService struct{
	authRepository  AuthRepository
}

func NewAuthService(authRepository AuthRepository)*AuthService{
return &AuthService{
	authRepository: authRepository,
}
}

func (s *AuthService) dCreateUser(ctx context.Context,payload *RegisterRequest) (*AuthResponce,error){

now := time.Now()



}