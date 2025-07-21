package storage

import (
	"context"

	"github.com/MukizuL/vk-test/internal/dto"
	filters2 "github.com/MukizuL/vk-test/internal/filters"
	"github.com/MukizuL/vk-test/internal/models"
	"github.com/MukizuL/vk-test/internal/storage/pg"
	"go.uber.org/fx"
)

//go:generate mockgen -source=storage.go -destination=mocks/storage.go -package=mockstorage

type Repo interface {
	CreateNewUser(ctx context.Context, login, passwordHash string) (string, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	CreateAd(ctx context.Context, login string, req *dto.CreateAdRequest) (*models.Ad, error)
	GetAds(ctx context.Context, filters filters2.Filters) ([]models.Ad, filters2.Metadata, error)
}

func newRepo(storage *pg.Storage) Repo {
	return storage
}

func Provide() fx.Option {
	return fx.Provide(newRepo)
}
