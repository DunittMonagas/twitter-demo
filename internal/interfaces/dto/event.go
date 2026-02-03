package dto

import "time"

// EventType represents the type of event being published
type EventType string

const (
	// TweetCreatedEvent is published when a new tweet is created
	TweetCreatedEvent EventType = "tweet.created"

	// TweetUpdatedEvent is published when a tweet is updated
	TweetUpdatedEvent EventType = "tweet.updated"

	// TweetDeletedEvent is published when a tweet is deleted
	TweetDeletedEvent EventType = "tweet.deleted"

	// FollowEvent is published when a user follows another user
	FollowEvent EventType = "follow.created"

	// UnfollowEvent is published when a user unfollows another user
	UnfollowEvent EventType = "follow.deleted"
)

// Event is a generic event wrapper for all domain events
type Event struct {
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// TweetCreatedEventData contains the data for a tweet.created event
// This is used for the Fan-Out pattern to distribute tweets to followers' timelines
type TweetCreatedEventData struct {
	TweetID   int64     `json:"tweet_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// TweetUpdatedEventData contains the data for a tweet.updated event
type TweetUpdatedEventData struct {
	TweetID   int64     `json:"tweet_id"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TweetDeletedEventData contains the data for a tweet.deleted event
type TweetDeletedEventData struct {
	TweetID int64 `json:"tweet_id"`
}

// FollowEventData contains the data for follow/unfollow events
type FollowEventData struct {
	FollowerID int64 `json:"follower_id"`
	FollowedID int64 `json:"followed_id"`
}

// NewEvent creates a new Event with the current timestamp
func NewEvent(eventType EventType, data interface{}) Event {
	return Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}
}
