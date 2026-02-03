package internal

import (
	"twitter-demo/internal/config"
	"twitter-demo/internal/infrastructure/repository"
	"twitter-demo/internal/interfaces/controller"
	"twitter-demo/internal/usecase"
	"twitter-demo/pkg"
)

type Container struct {
	UserController     controller.UserController
	TweetController    controller.TweetController
	FollowerController controller.FollowerController
	TimelineController controller.TimelineController
}

func NewContainer() (*Container, error) {

	db, err := pkg.NewPostgres(config.NewPostgresConfig())
	if err != nil {
		return nil, err
	}

	// Initialize Redis cache
	cache := pkg.NewRedisCache(config.NewRedisConfig())

	userRepository := repository.NewUser(db)
	userUsecase := usecase.NewUser(userRepository)
	userController := controller.NewUser(userUsecase)

	tweetRepository := repository.NewTweet(db)
	tweetUsecase := usecase.NewTweet(tweetRepository, userRepository)
	tweetController := controller.NewTweet(tweetUsecase)

	followerRepository := repository.NewFollower(db)
	followerUsecase := usecase.NewFollower(followerRepository, userRepository)
	followerController := controller.NewFollower(followerUsecase)

	timelineUsecase := usecase.NewTimeline(tweetRepository, cache)
	timelineController := controller.NewTimeline(timelineUsecase)

	return &Container{
		UserController:     userController,
		TweetController:    tweetController,
		FollowerController: followerController,
		TimelineController: timelineController,
	}, nil

}
