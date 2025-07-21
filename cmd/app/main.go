package main

import (
	"net/http"

	"github.com/MukizuL/vk-test/internal/config"
	"github.com/MukizuL/vk-test/internal/controller"
	mw "github.com/MukizuL/vk-test/internal/middleware"
	"github.com/MukizuL/vk-test/internal/migration"
	"github.com/MukizuL/vk-test/internal/router"
	"github.com/MukizuL/vk-test/internal/server"
	"github.com/MukizuL/vk-test/internal/services"
	"github.com/MukizuL/vk-test/internal/storage"
	"github.com/MukizuL/vk-test/internal/storage/pg"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		createApp(),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

func createApp() fx.Option {
	return fx.Options(
		config.Provide(),
		fx.Provide(zap.NewDevelopment),

		controller.Provide(),
		router.Provide(),
		server.Provide(),
		services.Provide(),
		mw.Provide(),
		migration.Provide(),

		pg.Provide(),
		storage.Provide(),
	)
}
