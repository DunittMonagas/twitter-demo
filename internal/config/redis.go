package config

import (
	"log"
	"os"
	"strconv"
)

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

func NewRedisConfig() RedisConfig {

	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Failed to convert REDIS_DB to int: %v", err)
	}

	return RedisConfig{
		Address:  os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}
}
