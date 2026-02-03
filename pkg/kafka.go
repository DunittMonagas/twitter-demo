package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"twitter-demo/internal/config"

	"github.com/IBM/sarama"
)

// MessageHandler is a function type for processing consumed messages.
// It receives the message key, value, and should return an error if processing fails.
type MessageHandler func(ctx context.Context, key, value []byte) error

// Producer defines the interface for publishing messages to Kafka.
// External code should depend on this interface, not on the concrete implementation.
type Producer interface {
	Publish(ctx context.Context, topic string, key string, message interface{}) error
	Close() error
}

// Consumer defines the interface for consuming messages from Kafka.
// External code should depend on this interface, not on the concrete implementation.
type Consumer interface {
	Consume(ctx context.Context, topics []string, handler MessageHandler) error
	Close() error
}

// kafkaProducer is the concrete implementation of Producer using sarama.
type kafkaProducer struct {
	producer sarama.SyncProducer
}

// kafkaConsumer is the concrete implementation of Consumer using sarama.
type kafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	handler       *consumerGroupHandler
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	handler MessageHandler
}

// NewKafkaProducer creates a new Kafka producer instance.
func NewKafkaProducer(cfg config.KafkaConfig) (Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 3

	producer, err := sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return &kafkaProducer{
		producer: producer,
	}, nil
}

// NewKafkaConsumer creates a new Kafka consumer instance.
// Note: The consumer uses the consumer group specified in the config.
func NewKafkaConsumer(cfg config.KafkaConfig) (Consumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_8_0_0
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer group: %w", err)
	}

	return &kafkaConsumer{
		consumerGroup: consumerGroup,
		handler:       &consumerGroupHandler{},
	}, nil
}

// Publish sends a message to the specified Kafka topic.
// The message is automatically serialized to JSON.
func (p *kafkaProducer) Publish(ctx context.Context, topic string, key string, message interface{}) error {
	// Serialize message to JSON
	valueBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Create Kafka message
	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(valueBytes),
	}

	// Send message
	partition, offset, err := p.producer.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to send message to topic %s: %w", topic, err)
	}

	log.Printf("Message sent to topic %s, partition %d, offset %d", topic, partition, offset)
	return nil
}

// Close closes the producer connection.
func (p *kafkaProducer) Close() error {
	return p.producer.Close()
}

// Consume starts consuming messages from the specified topics.
// This is a blocking operation that will continue until the context is canceled
// or an unrecoverable error occurs.
func (c *kafkaConsumer) Consume(ctx context.Context, topics []string, handler MessageHandler) error {
	c.handler.handler = handler

	for {
		// Check if context is done
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Consume messages
		if err := c.consumerGroup.Consume(ctx, topics, c.handler); err != nil {
			return fmt.Errorf("error from consumer: %w", err)
		}

		// Check if context was cancelled
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// Close closes the consumer connection.
func (c *kafkaConsumer) Close() error {
	return c.consumerGroup.Close()
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from a specific partition
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			// Process message with handler
			if err := h.handler(session.Context(), message.Key, message.Value); err != nil {
				log.Printf("Error processing message: %v", err)
				// Continue processing even if handler fails
				// You might want to implement retry logic or dead letter queue here
			} else {
				// Mark message as consumed on success
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			return nil
		}
	}
}
