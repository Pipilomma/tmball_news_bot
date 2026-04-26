package postgres

import (
	"context"
	"fmt"

	"tmballNews/internal/config"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func New(ctx context.Context, cfg *config.DBConfig) (*postgres, error) {
	postgresURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	connConfig, _ := pgx.ParseConfig(postgresURL)
	connStr := stdlib.RegisterConnConfig(connConfig)

	db, err := sqlx.Open("pgx", connStr)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	p := &postgres{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		cfg:     cfg,
	}

	return p, nil
}

func (p *postgres) Close(_ context.Context) error {
	return p.db.Close()
}
