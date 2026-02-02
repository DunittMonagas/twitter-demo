package main

import (
	"log"
	"twitter-demo/internal"

	"github.com/gin-gonic/gin"
)

func initRouter(c *internal.Container) *gin.Engine {

	router := gin.Default()
	apiV1 := router.Group("/api/v1")

	apiV1.POST("/users", c.UserController.CreateUser)
	apiV1.PUT("/users/:id", c.UserController.UpdateUser)

	apiV1.POST("/tweets", c.TweetController.CreateTweet)
	apiV1.PUT("/tweets/:id", c.TweetController.UpdateTweetByID)

	return router
}

func main() {

	container, err := internal.NewContainer()
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	router := initRouter(container)
	err = router.Run(":8081")
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
