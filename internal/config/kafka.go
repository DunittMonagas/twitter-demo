package config

import (
	"os"
	"strings"
)

type KafkaConfig struct {
	Brokers []string
	GroupID string
}

func NewKafkaConfig() KafkaConfig {
	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		brokersStr = "localhost:9092"
	}

	brokers := strings.Split(brokersStr, ",")

	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		groupID = "twitter-demo-default"
	}

	return KafkaConfig{
		Brokers: brokers,
		GroupID: groupID,
	}
}
