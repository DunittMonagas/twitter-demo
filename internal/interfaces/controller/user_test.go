package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"twitter-demo/internal/domain"
	"twitter-demo/internal/interfaces/dto"
	"twitter-demo/internal/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestUserController_GetUserByID_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	controller := NewUser(mockUsecase)

	expectedUser := domain.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUsecase.EXPECT().
		GetUserByID(gomock.Any(), int64(1)).
		Return(expectedUser, nil).
		Times(1)

	router := setupTestRouter()
	router.GET("/users/:id", controller.GetUserByID)

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.Username, response.Username)
	assert.Equal(t, expectedUser.Email, response.Email)
}

func TestUserController_GetUserByID_InvalidID(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	controller := NewUser(mockUsecase)

	router := setupTestRouter()
	router.GET("/users/:id", controller.GetUserByID)

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/users/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid id", response["error"])
}

func TestUserController_GetUserByID_UsecaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	controller := NewUser(mockUsecase)

	mockUsecase.EXPECT().
		GetUserByID(gomock.Any(), int64(1)).
		Return(domain.User{}, fmt.Errorf("database error")).
		Times(1)

	router := setupTestRouter()
	router.GET("/users/:id", controller.GetUserByID)

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "database error", response["error"])
}

func TestUserController_CreateUser_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	controller := NewUser(mockUsecase)

	createRequest := dto.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
	}

	expectedUser := domain.User{
		ID:        1,
		Username:  createRequest.Username,
		Email:     createRequest.Email,
		Password:  createRequest.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUsecase.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(expectedUser, nil).
		Times(1)

	router := setupTestRouter()
	router.POST("/users", controller.CreateUser)

	// Act
	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.Username, response.Username)
	assert.Equal(t, expectedUser.Email, response.Email)
}

func TestUserController_CreateUser_InvalidJSON(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	controller := NewUser(mockUsecase)

	router := setupTestRouter()
	router.POST("/users", controller.CreateUser)

	// Act
	invalidJSON := []byte(`{"username": "test", "email": }`)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
}

func TestUserController_CreateUser_UsecaseError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	controller := NewUser(mockUsecase)

	createRequest := dto.CreateUserRequest{
		Username: "duplicateuser",
		Email:    "duplicate@example.com",
		Password: "password123",
	}

	mockUsecase.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(domain.User{}, fmt.Errorf("email already exists")).
		Times(1)

	router := setupTestRouter()
	router.POST("/users", controller.CreateUser)

	// Act
	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "email already exists", response["error"])
}

func TestUserController_UpdateUser_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	controller := NewUser(mockUsecase)

	updateRequest := dto.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
		Password: "newpassword",
	}

	expectedUser := domain.User{
		ID:        1,
		Username:  updateRequest.Username,
		Email:     updateRequest.Email,
		Password:  updateRequest.Password,
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	mockUsecase.EXPECT().
		UpdateUser(gomock.Any(), int64(1), gomock.Any()).
		Return(expectedUser, nil).
		Times(1)

	router := setupTestRouter()
	router.PUT("/users/:id", controller.UpdateUser)

	// Act
	body, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.Username, response.Username)
	assert.Equal(t, expectedUser.Email, response.Email)
}

func TestUserController_GetAllUsers_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	controller := NewUser(mockUsecase)

	expectedUsers := []domain.User{
		{
			ID:        1,
			Username:  "user1",
			Email:     "user1@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Username:  "user2",
			Email:     "user2@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockUsecase.EXPECT().
		GetAllUsers(gomock.Any()).
		Return(expectedUsers, nil).
		Times(1)

	router := setupTestRouter()
	router.GET("/users", controller.GetAllUsers)

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response))
	assert.Equal(t, expectedUsers[0].ID, response[0].ID)
	assert.Equal(t, expectedUsers[1].ID, response[1].ID)
}
