package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"twitter-demo/internal/config"
	"twitter-demo/internal/interfaces/dto"
	"twitter-demo/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TimelineController interface {
	GetTimeline(ctx *gin.Context)
	HandleTweetCreated(ctx context.Context, key, value []byte) error
}

type Timeline struct {
	timelineUsecase usecase.TimelineUsecase
}

func NewTimeline(timelineUsecase usecase.TimelineUsecase) Timeline {
	return Timeline{
		timelineUsecase: timelineUsecase,
	}
}

func (t Timeline) GetTimeline(ctx *gin.Context) {
	// Get user ID from URL parameter
	userIDString := ctx.Param("id")
	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Get pagination parameters from query string
	var request dto.TimelineRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values if not provided
	if request.Limit <= 0 {
		request.Limit = config.DefaultLimit
	}

	if request.Limit > config.MaxLimit {
		request.Limit = config.MaxLimit
	}

	if request.Offset < 0 {
		request.Offset = 0
	}

	// Get timeline tweets
	tweets, err := t.timelineUsecase.GetTimeline(ctx, userID, request.Limit, request.Offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response
	response := dto.ToTimelineResponse(tweets, request.Limit, request.Offset)
	ctx.JSON(http.StatusOK, response)
}

// HandleTweetCreated is the Kafka message handler for tweet.created events.
// It implements the Fan-Out pattern by distributing the tweet to all followers' timelines.
func (t Timeline) HandleTweetCreated(ctx context.Context, key, value []byte) error {
	log.Printf("Received TweetCreatedEvent - Key: %s", string(key))

	// Deserialize the event
	var event dto.Event
	if err := json.Unmarshal(value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// Verify event type
	if event.Type != dto.TweetCreatedEvent {
		log.Printf("Unexpected event type: %s (expected: %s)", event.Type, dto.TweetCreatedEvent)
		return nil
	}

	// Deserialize the specific event data
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	var tweetData dto.TweetCreatedEventData
	if err := json.Unmarshal(dataBytes, &tweetData); err != nil {
		return fmt.Errorf("failed to parse tweet data: %w", err)
	}

	log.Printf("Tweet ID: %d, Author: %d", tweetData.TweetID, tweetData.UserID)
	log.Printf("   Content: %s", tweetData.Content)

	// Fan-Out: Distribute tweet to all followers' timelines
	if err := t.timelineUsecase.FanOutTweet(ctx, tweetData.UserID, tweetData.TweetID); err != nil {
		log.Printf("Fan-Out failed: %v", err)
		return err
	}

	log.Println("Fan-Out completed successfully")
	return nil
}
