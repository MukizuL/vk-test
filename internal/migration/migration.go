package migration

import (
	"database/sql"
	"embed"

	"github.com/MukizuL/vk-test/internal/config"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed "migrations/*.sql"
var embedMigrations embed.FS

type Migrator struct{}

func newMigrator(cfg *config.Config) (*Migrator, error) {
	db, err := sql.Open("pgx", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}

	if cfg.Dev {
		goose.Reset(db, "migrations")
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return nil, err
	}

	return &Migrator{}, nil
}

func Provide() fx.Option {
	return fx.Provide(newMigrator)
}
