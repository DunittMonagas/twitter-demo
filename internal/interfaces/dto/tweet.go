package dto

import (
	"time"
	"twitter-demo/internal/domain"
)

type CreateTweetRequest struct {
	UserID  int64  `json:"user_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateTweetRequest struct {
	Content string `json:"content" binding:"required"`
}

type TweetResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToTweetResponse(tweet domain.Tweet) TweetResponse {
	return TweetResponse{
		ID:        tweet.ID,
		UserID:    tweet.UserID,
		Content:   tweet.Content,
		CreatedAt: tweet.CreatedAt,
		UpdatedAt: tweet.UpdatedAt,
	}
}

func ToTweetDomain(request CreateTweetRequest) domain.Tweet {
	return domain.Tweet{
		UserID:  request.UserID,
		Content: request.Content,
	}
}

func ToUpdateTweetDomain(request UpdateTweetRequest) domain.Tweet {
	return domain.Tweet{
		Content: request.Content,
	}
}

type TimelineRequest struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type TimelineResponse struct {
	Tweets []TweetResponse `json:"tweets"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
	Total  int             `json:"total"`
}

func ToTimelineResponse(tweets []domain.Tweet, limit, offset int) TimelineResponse {
	tweetResponses := make([]TweetResponse, 0, len(tweets))
	for _, tweet := range tweets {
		tweetResponses = append(tweetResponses, ToTweetResponse(tweet))
	}

	return TimelineResponse{
		Tweets: tweetResponses,
		Limit:  limit,
		Offset: offset,
		Total:  len(tweetResponses),
	}
}
