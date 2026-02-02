package controller

import (
	"net/http"
	"twitter-demo/internal/interfaces/dto"
	"twitter-demo/internal/usecase"

	"github.com/gin-gonic/gin"
)

type FollowerController interface {
	FollowUser(ctx *gin.Context)
	UnfollowUser(ctx *gin.Context)
}

type Follower struct {
	followerUsecase usecase.FollowerUsecase
}

func NewFollower(followerUsecase usecase.FollowerUsecase) Follower {
	return Follower{
		followerUsecase: followerUsecase,
	}
}

func (f Follower) FollowUser(ctx *gin.Context) {

	followRequest := dto.FollowRequest{}
	if err := ctx.ShouldBindJSON(&followRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	follower, err := f.followerUsecase.FollowUser(ctx, followRequest.FollowerID, followRequest.FollowedID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToFollowerResponse(follower))
}

func (f Follower) UnfollowUser(ctx *gin.Context) {

	unfollowRequest := dto.UnfollowRequest{}
	if err := ctx.ShouldBindJSON(&unfollowRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := f.followerUsecase.UnfollowUser(ctx, unfollowRequest.FollowerID, unfollowRequest.FollowedID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully unfollowed user"})
}
