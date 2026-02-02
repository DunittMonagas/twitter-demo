package dto

import (
	"time"
	"twitter-demo/internal/domain"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserDomain(request CreateUserRequest) domain.User {
	return domain.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}
}

func ToUpdateUserDomain(request UpdateUserRequest) domain.User {
	return domain.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}
}
