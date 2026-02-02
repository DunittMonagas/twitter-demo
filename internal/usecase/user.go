package usecase

import (
	"context"
	"fmt"
	"twitter-demo/internal/domain"
	"twitter-demo/internal/infrastructure/repository"
)

type UserUsecase interface {
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	GetUserByID(ctx context.Context, id int64) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	UpdateUser(ctx context.Context, id int, user domain.User) (domain.User, error)
}

type User struct {
	userRepository repository.UserRepository
}

func NewUser(userRepository repository.UserRepository) User {
	return User{
		userRepository: userRepository,
	}
}

func (u User) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return u.userRepository.SelectAll(ctx)
}

func (u User) GetUserByID(ctx context.Context, id int64) (domain.User, error) {
	return u.userRepository.SelectByID(ctx, id)
}

func (u User) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return u.userRepository.SelectByEmail(ctx, email)
}

func (u User) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	return u.userRepository.SelectByUsername(ctx, username)
}

func (u User) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	return u.userRepository.Insert(ctx, user)
}

func (u User) UpdateUser(ctx context.Context, id int, user domain.User) (domain.User, error) {

	// Check if email already exists
	existingUserByEmail, err := u.userRepository.SelectByEmail(ctx, user.Email)
	if err != nil {
		return domain.User{}, err
	}
	if existingUserByEmail.ID != 0 {
		return domain.User{}, fmt.Errorf("email already exists")
	}

	// Check if username already exists
	existingUserByUsername, err := u.userRepository.SelectByUsername(ctx, user.Username)
	if err != nil {
		return domain.User{}, err
	}
	if existingUserByUsername.ID != 0 {
		return domain.User{}, fmt.Errorf("username already exists")
	}

	// Check if user exists
	existingUser, err := u.userRepository.SelectByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	if existingUser.ID == 0 {
		return domain.User{}, fmt.Errorf("user not found")
	}

	// Update user
	existingUser.Email = user.Email
	existingUser.Username = user.Username

	updatedUser, err := u.userRepository.UpdateByID(ctx, id, existingUser)
	if err != nil {
		return domain.User{}, err
	}

	return updatedUser, nil
}
