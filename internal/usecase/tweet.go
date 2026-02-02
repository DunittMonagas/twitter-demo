package usecase

import (
	"context"
	"fmt"
	"twitter-demo/internal/domain"
	"twitter-demo/internal/infrastructure/repository"
)

type TweetUsecase interface {
	GetTweetByID(ctx context.Context, id int64) (domain.Tweet, error)
	CreateTweet(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error)
	UpdateTweetByID(ctx context.Context, id int64, tweet domain.Tweet) (domain.Tweet, error)
}

type Tweet struct {
	tweetRepository repository.TweetRepository
	userRepository  repository.UserRepository
}

func NewTweet(tweetRepository repository.TweetRepository, userRepository repository.UserRepository) Tweet {
	return Tweet{
		tweetRepository: tweetRepository,
		userRepository:  userRepository,
	}
}

func (t Tweet) GetTweetByID(ctx context.Context, id int64) (domain.Tweet, error) {
	return t.tweetRepository.SelectByID(ctx, id)
}

func (t Tweet) CreateTweet(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error) {

	// Validate tweet
	err := t.validateTweet(ctx, tweet)
	if err != nil {
		fmt.Println("CreateTweet Error")
		fmt.Println(err)
		return domain.Tweet{}, err
	}

	newTweet, err := t.tweetRepository.Insert(ctx, tweet)
	if err != nil {
		fmt.Println("CreateTweet Error")
		fmt.Println(err)
		return domain.Tweet{}, err
	}

	fmt.Println("CreateTweet")
	fmt.Println(newTweet)

	return newTweet, nil
}

func (t Tweet) UpdateTweetByID(ctx context.Context, id int64, tweet domain.Tweet) (domain.Tweet, error) {

	// Validate tweet content
	if err := t.validateTweetContent(tweet.Content); err != nil {
		fmt.Println("UpdateTweetByID Error")
		fmt.Println(err)
		return domain.Tweet{}, err
	}

	// Check if tweet exists
	existingTweet, err := t.tweetRepository.SelectByID(ctx, id)
	if err != nil {
		fmt.Println("UpdateTweetByID Error")
		fmt.Println(err)
		return domain.Tweet{}, err
	}

	if existingTweet.ID == 0 {
		fmt.Println("UpdateTweetByID Error")
		fmt.Println(err)
		return domain.Tweet{}, fmt.Errorf("tweet not found")
	}

	// Update tweet content
	existingTweet.Content = tweet.Content

	updatedTweet, err := t.tweetRepository.UpdateByID(ctx, id, existingTweet)
	if err != nil {
		fmt.Println("UpdateTweetByID Error")
		fmt.Println(err)
		return domain.Tweet{}, err
	}

	return updatedTweet, nil
}

func (t Tweet) validateTweet(ctx context.Context, tweet domain.Tweet) error {

	// Validate content
	if err := t.validateTweetContent(tweet.Content); err != nil {
		return err
	}

	// Check if user exists
	existingUser, err := t.userRepository.SelectByID(ctx, tweet.UserID)
	if err != nil {
		fmt.Println("validateTweetByID Error SelectByID")
		fmt.Println(err)
		return err
	}
	if existingUser.ID == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (t Tweet) validateTweetContent(content string) error {

	// Check if content is empty
	if content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	// Check if content exceeds 280 characters
	if len(content) > 280 {
		return fmt.Errorf("content cannot exceed 280 characters")
	}

	return nil
}
