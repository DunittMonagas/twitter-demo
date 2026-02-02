package domain

import "time"

type Follower struct {
	ID         int64
	FollowerID int64
	FollowedID int64
	CreatedAt  time.Time
}
