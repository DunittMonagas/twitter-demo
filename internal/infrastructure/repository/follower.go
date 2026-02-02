package repository

import (
	"context"
	"database/sql"
	"fmt"
	"twitter-demo/internal/domain"
	"twitter-demo/pkg"
)

type FollowerRepository interface {
	Insert(ctx context.Context, follower domain.Follower) (domain.Follower, error)
	Delete(ctx context.Context, followerID, followedID int64) error
	SelectByFollowerAndFollowed(ctx context.Context, followerID, followedID int64) (domain.Follower, error)
}

type Follower struct {
	db *pkg.Postgres
}

func NewFollower(db *pkg.Postgres) Follower {
	return Follower{
		db: db,
	}
}

func (f Follower) Insert(ctx context.Context, follower domain.Follower) (domain.Follower, error) {

	var newFollower domain.Follower

	row := f.db.QueryRowContext(ctx,
		"INSERT INTO followers (follower_id, followed_id) VALUES ($1, $2) RETURNING id, follower_id, followed_id, created_at",
		follower.FollowerID, follower.FollowedID)

	err := row.Scan(&newFollower.ID, &newFollower.FollowerID, &newFollower.FollowedID, &newFollower.CreatedAt)
	if err != nil {
		fmt.Println("Insert Follower Error")
		fmt.Println(err)
		return domain.Follower{}, err
	}

	return newFollower, nil
}

func (f Follower) Delete(ctx context.Context, followerID, followedID int64) error {

	result, err := f.db.ExecContext(ctx,
		"DELETE FROM followers WHERE follower_id = $1 AND followed_id = $2",
		followerID, followedID)

	if err != nil {
		fmt.Println("Delete Follower Error")
		fmt.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("follower relationship not found")
	}

	return nil
}

func (f Follower) SelectByFollowerAndFollowed(ctx context.Context, followerID, followedID int64) (domain.Follower, error) {

	var follower domain.Follower

	row := f.db.QueryRowContext(ctx,
		"SELECT id, follower_id, followed_id, created_at FROM followers WHERE follower_id = $1 AND followed_id = $2",
		followerID, followedID)

	err := row.Scan(&follower.ID, &follower.FollowerID, &follower.FollowedID, &follower.CreatedAt)
	if err != nil {
		// If no rows found, return empty follower (ID will be 0) without error
		if err == sql.ErrNoRows {
			return domain.Follower{}, nil
		}
		return domain.Follower{}, err
	}

	return follower, nil
}
