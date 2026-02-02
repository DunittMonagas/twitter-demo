package internal

import (
	"twitter-demo/internal/config"
	"twitter-demo/internal/infrastructure/repository"
	"twitter-demo/internal/interfaces/controller"
	"twitter-demo/internal/usecase"
	"twitter-demo/pkg"
)

type Container struct {
	UserController  controller.UserController
	TweetController controller.TweetController
}

func NewContainer() (*Container, error) {

	db, err := pkg.NewPostgres(config.NewPostgresConfig())
	if err != nil {
		return nil, err
	}

	userRepository := repository.NewUser(db)
	userUsecase := usecase.NewUser(userRepository)
	userController := controller.NewUser(userUsecase)

	tweetRepository := repository.NewTweet(db)
	tweetUsecase := usecase.NewTweet(tweetRepository, userRepository)
	tweetController := controller.NewTweet(tweetUsecase)

	return &Container{
		UserController:  userController,
		TweetController: tweetController,
	}, nil

}
