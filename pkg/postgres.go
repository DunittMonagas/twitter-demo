package pkg

import (
	"database/sql"
	"fmt"
	"log"
	"twitter-demo/internal/config"

	_ "github.com/lib/pq"
)

type Postgres struct {
	*sql.DB
}

func NewPostgres(config config.PostgresConfig) (*Postgres, error) {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Database)

	db, err := sql.Open(config.Driver, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping Postgres: %v", err)
		return nil, err
	}

	return &Postgres{DB: db}, nil
}
