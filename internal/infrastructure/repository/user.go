package repository

import (
	"context"
	"database/sql"
	"fmt"
	"twitter-demo/internal/domain"
	"twitter-demo/pkg"
)

type UserRepository interface {
	SelectAll(ctx context.Context) ([]domain.User, error)
	SelectByID(ctx context.Context, id int64) (domain.User, error)
	SelectByEmail(ctx context.Context, email string) (domain.User, error)
	SelectByUsername(ctx context.Context, username string) (domain.User, error)
	Insert(ctx context.Context, user domain.User) (domain.User, error)
	UpdateByID(ctx context.Context, id int64, user domain.User) (domain.User, error)
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
		fmt.Println("SelectAll Error")
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	fmt.Println("SelectAll Success")
	fmt.Println(rows)

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	fmt.Println(users)

	return users, nil
}

func (u User) SelectByID(ctx context.Context, id int64) (domain.User, error) {

	var user domain.User

	row := u.db.QueryRowContext(ctx, "SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = $1", id)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		// If no rows found, return empty user (ID will be 0) without error
		if err == sql.ErrNoRows {
			return domain.User{}, nil
		}
		return domain.User{}, err
	}

	return user, nil
}

func (u User) SelectByEmail(ctx context.Context, email string) (domain.User, error) {

	var user domain.User

	row := u.db.QueryRowContext(ctx, "SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = $1", email)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		// If no rows found, return empty user (ID will be 0) without error
		if err == sql.ErrNoRows {
			return domain.User{}, nil
		}
		return domain.User{}, err
	}

	return user, nil
}

func (u User) SelectByUsername(ctx context.Context, username string) (domain.User, error) {

	var user domain.User

	row := u.db.QueryRowContext(ctx, "SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = $1", username)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		// If no rows found, return empty user (ID will be 0) without error
		if err == sql.ErrNoRows {
			return domain.User{}, nil
		}
		return domain.User{}, err
	}

	return user, nil
}

func (u User) Insert(ctx context.Context, user domain.User) (domain.User, error) {

	var newUser domain.User

	row := u.db.QueryRowContext(ctx, "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, username, email, password, created_at, updated_at", user.Username, user.Email, user.Password)

	err := row.Scan(&newUser.ID, &newUser.Username, &newUser.Email, &newUser.Password, &newUser.CreatedAt, &newUser.UpdatedAt)
	if err != nil {
		fmt.Println("Insert Error")
		fmt.Println(err)
		return domain.User{}, err
	}

	return newUser, nil
}

func (u User) UpdateByID(ctx context.Context, id int64, user domain.User) (domain.User, error) {

	var updatedUser domain.User

	row := u.db.QueryRowContext(ctx, "UPDATE users SET username = $1, email = $2, password = $3 WHERE id = $4 RETURNING id, username, email, password, created_at, updated_at", user.Username, user.Email, user.Password, id)

	err := row.Scan(&updatedUser.ID, &updatedUser.Username, &updatedUser.Email, &updatedUser.Password, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)
	if err != nil {
		fmt.Println("UpdateByID Error")
		fmt.Println(err)
		return domain.User{}, err
	}

	return updatedUser, nil
}
