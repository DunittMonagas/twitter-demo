package config

import (
	"log"
	"os"
	"strconv"
)

const (
	PostgresDriver = "postgres"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Driver   string
}

func NewPostgresConfig() PostgresConfig {

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Failed to convert DB_PORT to int: %v", err)
	}

	return PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		Driver:   PostgresDriver,
	}
}
