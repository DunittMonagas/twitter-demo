package usecase

import (
	"context"
	"twitter-demo/internal/config"
	"twitter-demo/internal/domain"
	"twitter-demo/internal/infrastructure/repository"
)

type TimelineUsecase interface {
	GetTimeline(ctx context.Context, userID int64, limit, offset int) ([]domain.Tweet, error)
}

type Timeline struct {
	tweetRepository repository.TweetRepository
}

func NewTimeline(tweetRepository repository.TweetRepository) Timeline {
	return Timeline{
		tweetRepository: tweetRepository,
	}
}

func (t Timeline) GetTimeline(ctx context.Context, userID int64, limit, offset int) ([]domain.Tweet, error) {
	// Set default and max values for pagination
	if limit <= 0 {
		limit = config.DefaultLimit
	}
	if limit > config.MaxLimit {
		limit = config.MaxLimit
	}
	if offset < 0 {
		offset = 0
	}

	tweets, err := t.tweetRepository.SelectTimelineTweets(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}
