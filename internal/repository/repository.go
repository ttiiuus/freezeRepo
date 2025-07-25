package repository

import (
	"auth/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPgxPool(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Server.Host, cfg.Postgres.Server.Port, cfg.Postgres.Name,
	)
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	//poolConfig.MaxConns = 10
	//poolConfig.MinConns = 2
	//poolConfig.HealthCheckPeriod = 0 // например, отключить

	ctx := context.Background()

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB pool: %w", err)
	}

	// small checkhealth
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return pool, nil
}
