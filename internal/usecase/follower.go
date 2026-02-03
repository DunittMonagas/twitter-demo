package usecase

import (
	"context"
	"fmt"
	"twitter-demo/internal/domain"
	"twitter-demo/internal/infrastructure/repository"
)

type FollowerUsecase interface {
	FollowUser(ctx context.Context, followerID, followedID int64) (domain.Follower, error)
	UnfollowUser(ctx context.Context, followerID, followedID int64) error
}

type Follower struct {
	followerRepository repository.FollowerRepository
	userRepository     repository.UserRepository
}

func NewFollower(followerRepository repository.FollowerRepository, userRepository repository.UserRepository) Follower {
	return Follower{
		followerRepository: followerRepository,
		userRepository:     userRepository,
	}
}

func (f Follower) FollowUser(ctx context.Context, followerID, followedID int64) (domain.Follower, error) {

	// Validate that follower ID and followed ID are different
	if followerID == followedID {
		return domain.Follower{}, fmt.Errorf("cannot follow yourself")
	}

	// Check if follower user exists
	followerUser, err := f.userRepository.SelectByID(ctx, followerID)
	if err != nil {
		return domain.Follower{}, err
	}
	if followerUser.ID == 0 {
		return domain.Follower{}, fmt.Errorf("follower user not found")
	}

	// Check if followed user exists
	followedUser, err := f.userRepository.SelectByID(ctx, followedID)
	if err != nil {
		return domain.Follower{}, err
	}
	if followedUser.ID == 0 {
		return domain.Follower{}, fmt.Errorf("followed user not found")
	}

	// Check if relationship already exists
	existingFollower, err := f.followerRepository.SelectByFollowerAndFollowed(ctx, followerID, followedID)
	if err != nil {
		return domain.Follower{}, err
	}
	if existingFollower.ID != 0 {
		return domain.Follower{}, fmt.Errorf("already following this user")
	}

	// Create follower relationship
	newFollower := domain.Follower{
		FollowerID: followerID,
		FollowedID: followedID,
	}

	createdFollower, err := f.followerRepository.Insert(ctx, newFollower)
	if err != nil {
		return domain.Follower{}, err
	}

	return createdFollower, nil
}

func (f Follower) UnfollowUser(ctx context.Context, followerID, followedID int64) error {

	// Validate that follower ID and followed ID are different
	if followerID == followedID {
		return fmt.Errorf("invalid unfollow operation")
	}

	// Check if follower user exists
	followerUser, err := f.userRepository.SelectByID(ctx, followerID)
	if err != nil {
		return err
	}
	if followerUser.ID == 0 {
		return fmt.Errorf("follower user not found")
	}

	// Check if followed user exists
	followedUser, err := f.userRepository.SelectByID(ctx, followedID)
	if err != nil {
		return err
	}
	if followedUser.ID == 0 {
		return fmt.Errorf("followed user not found")
	}

	// Check if relationship exists
	existingFollower, err := f.followerRepository.SelectByFollowerAndFollowed(ctx, followerID, followedID)
	if err != nil {
		return err
	}
	if existingFollower.ID == 0 {
		return fmt.Errorf("not following this user")
	}

	// Delete follower relationship
	err = f.followerRepository.Delete(ctx, followerID, followedID)
	if err != nil {
		return err
	}

	return nil
}
