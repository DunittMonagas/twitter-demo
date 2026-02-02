package dto

import (
	"time"
	"twitter-demo/internal/domain"
)

type FollowRequest struct {
	FollowerID int64 `json:"follower_id" binding:"required"`
	FollowedID int64 `json:"followed_id" binding:"required"`
}

type UnfollowRequest struct {
	FollowerID int64 `json:"follower_id" binding:"required"`
	FollowedID int64 `json:"followed_id" binding:"required"`
}

type FollowerResponse struct {
	ID         int64     `json:"id"`
	FollowerID int64     `json:"follower_id"`
	FollowedID int64     `json:"followed_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToFollowerResponse(follower domain.Follower) FollowerResponse {
	return FollowerResponse{
		ID:         follower.ID,
		FollowerID: follower.FollowerID,
		FollowedID: follower.FollowedID,
		CreatedAt:  follower.CreatedAt,
	}
}

func ToFollowerDomain(request FollowRequest) domain.Follower {
	return domain.Follower{
		FollowerID: request.FollowerID,
		FollowedID: request.FollowedID,
	}
}
