package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"twitter-demo/internal"
	"twitter-demo/internal/config"
)

func main() {
	log.Println("========================================")
	log.Println("KAFKA WORKER - FAN-OUT PROCESSOR")
	log.Println("========================================")

	// Initialize worker container
	container, err := internal.NewWorkerContainer()
	if err != nil {
		log.Fatalf("Failed to create worker container: %v", err)
	}
	defer func() {
		log.Println("Closing consumer...")
		if err := container.Consumer.Close(); err != nil {
			log.Printf("Error closing consumer: %v", err)
		}
	}()

	log.Println("Worker container initialized")

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for termination signals
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// Start consuming tweets topic
	topics := []string{config.TopicTweets}
	log.Printf("📡 Listening for events on topic: %s", config.TopicTweets)
	log.Println("Press Ctrl+C to stop...")
	log.Println("========================================")

	// Start consuming messages in a goroutine
	go func() {
		if err := container.Consumer.Consume(ctx, topics, container.TimelineController.HandleTweetCreated); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	// Wait for termination signal
	<-sigterm
	log.Println("\n========================================")
	log.Println("Shutting down worker gracefully...")
	log.Println("========================================")
}
