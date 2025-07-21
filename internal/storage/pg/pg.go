package pg

import (
	"context"
	"fmt"

	"github.com/MukizuL/vk-test/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type Storage struct {
	conn *pgxpool.Pool
}

func newStorage(lc fx.Lifecycle, cfg *config.Config) (*Storage, error) {
	dbpool, err := pgxpool.New(context.Background(), cfg.DB.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create db pool: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return dbpool.Ping(ctx)
		},
		OnStop: func(ctx context.Context) error {
			dbpool.Close()
			return nil
		},
	})

	return &Storage{
		conn: dbpool,
	}, nil
}

func Provide() fx.Option {
	return fx.Provide(newStorage)
}
