package controller

import (
	"net/http"
	"strconv"
	"twitter-demo/internal/config"
	"twitter-demo/internal/interfaces/dto"
	"twitter-demo/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TimelineController interface {
	GetTimeline(ctx *gin.Context)
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
