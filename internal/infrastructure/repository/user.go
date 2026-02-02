package repository

import (
	"context"
	"twitter-demo/internal/domain"
	"twitter-demo/pkg"
)

type UserRepository interface {
	SelectAll(ctx context.Context) ([]domain.User, error)
	SelectByID(ctx context.Context, id int) (domain.User, error)
	SelectByEmail(ctx context.Context, email string) (domain.User, error)
	SelectByUsername(ctx context.Context, username string) (domain.User, error)
	Insert(ctx context.Context, user domain.User) (domain.User, error)
	UpdateByID(ctx context.Context, id int, user domain.User) (domain.User, error)
}

type User struct {
	db *pkg.Postgres
}

func NewUser(db *pkg.Postgres) User {
	return User{
		db: db,
	}
}

func (u User) SelectAll(ctx context.Context) ([]domain.User, error) {
	rows, err := u.db.QueryContext(ctx, "SELECT id, username, email, password, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u User) SelectByID(ctx context.Context, id int) (domain.User, error) {
	row := u.db.QueryRowContext(ctx, "SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = $1", id)
	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u User) SelectByEmail(ctx context.Context, email string) (domain.User, error) {
	row := u.db.QueryRowContext(ctx, "SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = $1", email)
	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u User) SelectByUsername(ctx context.Context, username string) (domain.User, error) {
	row := u.db.QueryRowContext(ctx, "SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = $1", username)
	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u User) Insert(ctx context.Context, user domain.User) (domain.User, error) {
	row := u.db.QueryRowContext(ctx, "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, username, email, password, created_at, updated_at", user.Username, user.Email, user.Password)
	var newUser domain.User
	err := row.Scan(&newUser.ID, &newUser.Username, &newUser.Email, &newUser.Password, &newUser.CreatedAt, &newUser.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return newUser, nil
}

func (u User) UpdateByID(ctx context.Context, id int, user domain.User) (domain.User, error) {
	row := u.db.QueryRowContext(ctx, "UPDATE users SET username = $1, email = $2, password = $3 WHERE id = $4 RETURNING id, username, email, password, created_at, updated_at", user.Username, user.Email, user.Password, id)
	var updatedUser domain.User
	err := row.Scan(&updatedUser.ID, &updatedUser.Username, &updatedUser.Email, &updatedUser.Password, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return updatedUser, nil
}
