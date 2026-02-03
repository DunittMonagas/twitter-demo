package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"
	"twitter-demo/internal/domain"
	"twitter-demo/internal/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUser_CreateUser_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := NewUser(mockRepo)

	newUser := domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	expectedUser := domain.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock to validate that email does not exist
	mockRepo.EXPECT().
		SelectByEmail(gomock.Any(), newUser.Email).
		Return(domain.User{}, nil).
		Times(1)

	// Mock to validate that username does not exist
	mockRepo.EXPECT().
		SelectByUsername(gomock.Any(), newUser.Username).
		Return(domain.User{}, nil).
		Times(1)

	// Mock to insert the new user
	mockRepo.EXPECT().
		Insert(gomock.Any(), newUser).
		Return(expectedUser, nil).
		Times(1)

	// Act
	createdUser, err := usecase.CreateUser(context.Background(), newUser)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, createdUser.ID)
	assert.Equal(t, expectedUser.Username, createdUser.Username)
	assert.Equal(t, expectedUser.Email, createdUser.Email)
}

func TestUser_CreateUser_EmailAlreadyExists(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := NewUser(mockRepo)

	newUser := domain.User{
		Username: "testuser",
		Email:    "existing@example.com",
		Password: "hashedpassword",
	}

	existingUser := domain.User{
		ID:       1,
		Username: "existinguser",
		Email:    "existing@example.com",
	}

	// Mock to validate that email already exists
	mockRepo.EXPECT().
		SelectByEmail(gomock.Any(), newUser.Email).
		Return(existingUser, nil).
		Times(1)

	// Act
	_, err := usecase.CreateUser(context.Background(), newUser)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "email already exists", err.Error())
}

func TestUser_CreateUser_UsernameAlreadyExists(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := NewUser(mockRepo)

	newUser := domain.User{
		Username: "existinguser",
		Email:    "new@example.com",
		Password: "hashedpassword",
	}

	existingUser := domain.User{
		ID:       1,
		Username: "existinguser",
		Email:    "other@example.com",
	}

	// Mock to validate that email does not exist
	mockRepo.EXPECT().
		SelectByEmail(gomock.Any(), newUser.Email).
		Return(domain.User{}, nil).
		Times(1)

	// Mock to validate that username already exists
	mockRepo.EXPECT().
		SelectByUsername(gomock.Any(), newUser.Username).
		Return(existingUser, nil).
		Times(1)

	// Act
	_, err := usecase.CreateUser(context.Background(), newUser)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "username already exists", err.Error())
}

func TestUser_UpdateUser_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := NewUser(mockRepo)

	userID := int64(1)
	updateData := domain.User{
		Username: "updateduser",
		Email:    "updated@example.com",
		Password: "newpassword",
	}

	existingUser := domain.User{
		ID:        userID,
		Username:  "olduser",
		Email:     "old@example.com",
		Password:  "oldpassword",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	}

	updatedUser := domain.User{
		ID:        userID,
		Username:  "updateduser",
		Email:     "updated@example.com",
		Password:  "newpassword",
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: time.Now(),
	}

	// Mocks for validation
	mockRepo.EXPECT().
		SelectByEmail(gomock.Any(), updateData.Email).
		Return(domain.User{}, nil).
		Times(1)

	mockRepo.EXPECT().
		SelectByUsername(gomock.Any(), updateData.Username).
		Return(domain.User{}, nil).
		Times(1)

	// Mock to verify that user exists
	mockRepo.EXPECT().
		SelectByID(gomock.Any(), userID).
		Return(existingUser, nil).
		Times(1)

	// Mock to update the user
	mockRepo.EXPECT().
		UpdateByID(gomock.Any(), userID, gomock.Any()).
		Return(updatedUser, nil).
		Times(1)

	// Act
	result, err := usecase.UpdateUser(context.Background(), userID, updateData)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.ID, result.ID)
	assert.Equal(t, updatedUser.Username, result.Username)
	assert.Equal(t, updatedUser.Email, result.Email)
}

func TestUser_UpdateUser_UserNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := NewUser(mockRepo)

	userID := int64(999)
	updateData := domain.User{
		Username: "updateduser",
		Email:    "updated@example.com",
		Password: "newpassword",
	}

	// Mocks for validation
	mockRepo.EXPECT().
		SelectByEmail(gomock.Any(), updateData.Email).
		Return(domain.User{}, nil).
		Times(1)

	mockRepo.EXPECT().
		SelectByUsername(gomock.Any(), updateData.Username).
		Return(domain.User{}, nil).
		Times(1)

	// Mock to verify that user does not exist (ID = 0)
	mockRepo.EXPECT().
		SelectByID(gomock.Any(), userID).
		Return(domain.User{}, nil).
		Times(1)

	// Act
	_, err := usecase.UpdateUser(context.Background(), userID, updateData)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestUser_GetUserByID_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := NewUser(mockRepo)

	userID := int64(1)
	expectedUser := domain.User{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.EXPECT().
		SelectByID(gomock.Any(), userID).
		Return(expectedUser, nil).
		Times(1)

	// Act
	user, err := usecase.GetUserByID(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Email, user.Email)
}

func TestUser_GetUserByID_Error(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	usecase := NewUser(mockRepo)

	userID := int64(1)
	expectedError := fmt.Errorf("database error")

	mockRepo.EXPECT().
		SelectByID(gomock.Any(), userID).
		Return(domain.User{}, expectedError).
		Times(1)

	// Act
	_, err := usecase.GetUserByID(context.Background(), userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}
