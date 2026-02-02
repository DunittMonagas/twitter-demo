package repository

import (
	"context"
	"database/sql"
	"fmt"
	"twitter-demo/internal/domain"
	"twitter-demo/pkg"
)

type TweetRepository interface {
	SelectByID(ctx context.Context, id int64) (domain.Tweet, error)
	Insert(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error)
	UpdateByID(ctx context.Context, id int64, tweet domain.Tweet) (domain.Tweet, error)
}

type Tweet struct {
	db *pkg.Postgres
}

func NewTweet(db *pkg.Postgres) Tweet {
	return Tweet{
		db: db,
	}
}

func (t Tweet) SelectByID(ctx context.Context, id int64) (domain.Tweet, error) {

	var tweet domain.Tweet

	row := t.db.QueryRowContext(ctx, "SELECT id, user_id, content, created_at, updated_at FROM tweets WHERE id = $1", id)

	err := row.Scan(&tweet.ID, &tweet.UserID, &tweet.Content, &tweet.CreatedAt, &tweet.UpdatedAt)
	if err != nil {
		// If no rows found, return empty tweet (ID will be 0) without error
		if err == sql.ErrNoRows {
			return domain.Tweet{}, nil
		}
		return domain.Tweet{}, err
	}

	return tweet, nil
}

func (t Tweet) Insert(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error) {

	var newTweet domain.Tweet

	row := t.db.QueryRowContext(ctx, "INSERT INTO tweets (user_id, content) VALUES ($1, $2) RETURNING id, user_id, content, created_at, updated_at", tweet.UserID, tweet.Content)

	err := row.Scan(&newTweet.ID, &newTweet.UserID, &newTweet.Content, &newTweet.CreatedAt, &newTweet.UpdatedAt)
	if err != nil {
		fmt.Println("Insert Error")
		fmt.Println(err)
		return domain.Tweet{}, err
	}

	return newTweet, nil
}

func (t Tweet) UpdateByID(ctx context.Context, id int64, tweet domain.Tweet) (domain.Tweet, error) {

	var updatedTweet domain.Tweet

	row := t.db.QueryRowContext(ctx, "UPDATE tweets SET content = $1 WHERE id = $2 RETURNING id, user_id, content, created_at, updated_at", tweet.Content, id)

	err := row.Scan(&updatedTweet.ID, &updatedTweet.UserID, &updatedTweet.Content, &updatedTweet.CreatedAt, &updatedTweet.UpdatedAt)
	if err != nil {
		fmt.Println("UpdateByID Error")
		fmt.Println(err)
		return domain.Tweet{}, err
	}

	return updatedTweet, nil
}
