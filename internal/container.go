package internal

import (
	"log"
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

	// Initialize Kafka producer
	kafkaConfig := config.NewKafkaConfig()
	producer, err := pkg.NewKafkaProducer(kafkaConfig)
	if err != nil {
		log.Printf("Warning: Failed to create Kafka producer: %v", err)
		// Continue without Kafka - events won't be published but API will work
	}

	userRepository := repository.NewUser(db)
	userUsecase := usecase.NewUser(userRepository)
	userController := controller.NewUser(userUsecase)

	tweetRepository := repository.NewTweet(db)
	tweetUsecase := usecase.NewTweet(tweetRepository, userRepository, producer)
	tweetController := controller.NewTweet(tweetUsecase)

	followerRepository := repository.NewFollower(db)
	followerUsecase := usecase.NewFollower(followerRepository, userRepository)
	followerController := controller.NewFollower(followerUsecase)

	timelineUsecase := usecase.NewTimeline(tweetRepository, followerRepository, cache)
	timelineController := controller.NewTimeline(timelineUsecase)

	return &Container{
		UserController:     userController,
		TweetController:    tweetController,
		FollowerController: followerController,
		TimelineController: timelineController,
	}, nil

}

// WorkerContainer holds dependencies for the Kafka worker service
type WorkerContainer struct {
	TimelineController controller.TimelineController
	Consumer           pkg.Consumer
}

// NewWorkerContainer creates a new container for the worker service
func NewWorkerContainer() (*WorkerContainer, error) {

	db, err := pkg.NewPostgres(config.NewPostgresConfig())
	if err != nil {
		return nil, err
	}

	// Initialize Redis cache
	cache := pkg.NewRedisCache(config.NewRedisConfig())

	// Initialize Kafka consumer
	kafkaConfig := config.NewKafkaConfig()
	consumer, err := pkg.NewKafkaConsumer(kafkaConfig)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	tweetRepository := repository.NewTweet(db)
	followerRepository := repository.NewFollower(db)

	// Initialize use cases
	timelineUsecase := usecase.NewTimeline(tweetRepository, followerRepository, cache)

	// Initialize controllers
	timelineController := controller.NewTimeline(timelineUsecase)

	return &WorkerContainer{
		TimelineController: timelineController,
		Consumer:           consumer,
	}, nil
}
