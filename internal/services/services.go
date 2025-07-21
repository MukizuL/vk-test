package services

import (
	"github.com/MukizuL/vk-test/internal/config"
	"github.com/MukizuL/vk-test/internal/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Services struct {
	storage storage.Repo
	logger  *zap.Logger
	key     []byte
}

func newServices(storage storage.Repo, logger *zap.Logger, cfg *config.Config) *Services {
	return &Services{
		storage: storage,
		logger:  logger,
		key:     []byte(cfg.Secret),
	}
}

func Provide() fx.Option {
	return fx.Provide(newServices)
}
