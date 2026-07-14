package auth

import "github.com/google/uuid"

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=50"`
	LastName  string `json:"last_name" validate:"required,min=3,max=50"`
	Username  string `json:"username" validate:"required,username"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type AuthResponce struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	User         UserDTO   `json:"user"`
}

type UserDTO struct {
    ID        uuid.UUID `json:"id"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    Role      string    `json:"role"`
    ImageID   *uuid.UUID `json:"image_id,omitempty"`
}
