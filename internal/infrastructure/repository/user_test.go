package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"twitter-demo/internal/domain"
	"twitter-demo/pkg"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUser_SelectByID_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	postgres := &pkg.Postgres{DB: db}
	repo := NewUser(postgres)

	expectedUser := domain.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery("SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = \\$1").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	// Act
	user, err := repo.SelectByID(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_SelectByID_NotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	postgres := &pkg.Postgres{DB: db}
	repo := NewUser(postgres)

	mock.ExpectQuery("SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = \\$1").
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	// Act
	user, err := repo.SelectByID(context.Background(), 999)

	// Assert
	assert.NoError(t, err)             // Repository returns nil error when not found
	assert.Equal(t, int64(0), user.ID) // Empty user
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_SelectByEmail_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	postgres := &pkg.Postgres{DB: db}
	repo := NewUser(postgres)

	expectedUser := domain.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email, expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery("SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = \\$1").
		WithArgs("test@example.com").
		WillReturnRows(rows)

	// Act
	user, err := repo.SelectByEmail(context.Background(), "test@example.com")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_Insert_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	postgres := &pkg.Postgres{DB: db}
	repo := NewUser(postgres)

	newUser := domain.User{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "hashedpassword",
	}

	expectedTime := time.Now()
	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
		AddRow(int64(1), newUser.Username, newUser.Email, newUser.Password, expectedTime, expectedTime)

	mock.ExpectQuery("INSERT INTO users \\(username, email, password\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id, username, email, password, created_at, updated_at").
		WithArgs(newUser.Username, newUser.Email, newUser.Password).
		WillReturnRows(rows)

	// Act
	createdUser, err := repo.Insert(context.Background(), newUser)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), createdUser.ID)
	assert.Equal(t, newUser.Username, createdUser.Username)
	assert.Equal(t, newUser.Email, createdUser.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_Insert_Error(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	postgres := &pkg.Postgres{DB: db}
	repo := NewUser(postgres)

	newUser := domain.User{
		Username: "duplicateuser",
		Email:    "duplicate@example.com",
		Password: "hashedpassword",
	}

	mock.ExpectQuery("INSERT INTO users \\(username, email, password\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id, username, email, password, created_at, updated_at").
		WithArgs(newUser.Username, newUser.Email, newUser.Password).
		WillReturnError(sql.ErrConnDone)

	// Act
	_, err = repo.Insert(context.Background(), newUser)

	// Assert
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
