package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"twitter-demo/internal"
)

func initRouter(c *internal.Container) *gin.Engine {

	router := gin.Default()
	apiV1 := router.Group("/api/v1")

	apiV1.GET("/users", c.UserController.GetAllUsers)
	apiV1.GET("/users/:id", c.UserController.GetUserByID)

	return router

}

func main() {

	container, err := internal.NewContainer()
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	router := initRouter(container)
	err = router.Run()
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
