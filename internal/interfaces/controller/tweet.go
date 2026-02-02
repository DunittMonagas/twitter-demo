package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"twitter-demo/internal/interfaces/dto"
	"twitter-demo/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TweetController interface {
	GetTweetByID(ctx *gin.Context)
	CreateTweet(ctx *gin.Context)
	UpdateTweetByID(ctx *gin.Context)
}

type Tweet struct {
	tweetUsecase usecase.TweetUsecase
}

func NewTweet(tweetUsecase usecase.TweetUsecase) Tweet {
	return Tweet{
		tweetUsecase: tweetUsecase,
	}
}

func (t Tweet) GetTweetByID(ctx *gin.Context) {

	idString := ctx.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tweet, err := t.tweetUsecase.GetTweetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if tweet.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "tweet not found"})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToTweetResponse(tweet))
}

func (t Tweet) CreateTweet(ctx *gin.Context) {

	createTweetRequest := dto.CreateTweetRequest{}
	if err := ctx.ShouldBindJSON(&createTweetRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweet := dto.ToTweetDomain(createTweetRequest)
	newTweet, err := t.tweetUsecase.CreateTweet(ctx, tweet)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToTweetResponse(newTweet))
}

func (t Tweet) UpdateTweetByID(ctx *gin.Context) {

	updateTweetRequest := dto.UpdateTweetRequest{}
	if err := ctx.ShouldBindJSON(&updateTweetRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idString := ctx.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tweet := dto.ToUpdateTweetDomain(updateTweetRequest)
	updatedTweet, err := t.tweetUsecase.UpdateTweetByID(ctx, id, tweet)
	if err != nil {
		fmt.Println("UpdateTweetByID Error")
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToTweetResponse(updatedTweet))
}
