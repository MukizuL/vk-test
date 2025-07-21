package controller

import (
	"github.com/MukizuL/vk-test/internal/config"
	"github.com/MukizuL/vk-test/internal/services"
	"go.uber.org/fx"
)

type Controller struct {
	domain  string
	service *services.Services
}

func newController(service *services.Services, cfg *config.Config) *Controller {
	return &Controller{
		domain:  cfg.Domain,
		service: service,
	}
}

func Provide() fx.Option {
	return fx.Provide(newController)
}
