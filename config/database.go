package config

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres(cfg *Config) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.PGUser,
		cfg.PGPassword,
		cfg.PGHost,
		cfg.PGPort,
		cfg.PGDBName,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	return pool
}
