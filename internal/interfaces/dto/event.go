package dto

import "time"

// EventType represents the type of event being published
type EventType string

const (
	// TweetCreatedEvent is published when a new tweet is created
	TweetCreatedEvent EventType = "tweet.created"
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

// NewEvent creates a new Event with the current timestamp
func NewEvent(eventType EventType, data interface{}) Event {
	return Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}
}
