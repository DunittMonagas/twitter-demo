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

	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Fatalf("Failed to convert POSTGRES_PORT to int: %v", err)
	}

	return PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     port,
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
		Driver:   PostgresDriver,
	}
}
