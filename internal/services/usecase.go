package services

import (
	"context"
	"strings"

	"github.com/MukizuL/vk-test/internal/dto"
	"github.com/MukizuL/vk-test/internal/errs"
	filters2 "github.com/MukizuL/vk-test/internal/filters"
	"github.com/greatcloak/decimal"
	"go.uber.org/zap"
)

func (s *Services) CreateAd(ctx context.Context, userID string, req *dto.CreateAdRequest) (*dto.CreateAdResponse, error) {
	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user by id", zap.String("user_id", userID), zap.Error(err))
		return nil, errs.TransformPGErrors(err)
	}

	req.Price = req.Price.Mul(decimal.NewFromInt(100))

	ad, err := s.storage.CreateAd(ctx, user.Login, req)
	if err != nil {
		s.logger.Error("failed to create an ad", zap.String("login", user.Login), zap.Any("request", req), zap.Error(err))
		return nil, errs.TransformPGErrors(err)
	}

	response := &dto.CreateAdResponse{
		Login:       user.Login,
		Title:       ad.Title,
		Description: ad.Description,
		ImageURL:    ad.ImageURL,
		Price:       ad.Price.Div(decimal.NewFromInt(100)),
		CreatedAt:   ad.CreatedAt,
	}

	return response, nil
}

func (s *Services) GetAds(ctx context.Context, userID string, filters filters2.Filters) ([]dto.GetAdsResponse, filters2.Metadata, error) {
	if strings.Contains(filters.Sort, "date") {
		filters.Sort = strings.ReplaceAll(filters.Sort, "date", "created_at")
	}

	ads, metadata, err := s.storage.GetAds(ctx, filters)
	if err != nil {
		s.logger.Error("failed to get ads", zap.String("user_id", userID), zap.Error(err))
		return nil, filters2.Metadata{}, errs.TransformPGErrors(err)
	}

	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user by id", zap.String("user_id", userID), zap.Error(err))
		return nil, filters2.Metadata{}, errs.TransformPGErrors(err)
	}

	var response []dto.GetAdsResponse

	for _, ad := range ads {
		temp := dto.GetAdsResponse{
			Login:       ad.UserLogin,
			Owned:       false,
			Title:       ad.Title,
			Description: ad.Description,
			ImageURL:    ad.ImageURL,
			Price:       ad.Price.Div(decimal.NewFromInt(100)),
			CreatedAt:   ad.CreatedAt,
		}

		if ad.UserLogin == user.Login {
			temp.Owned = true
		}

		response = append(response, temp)
	}

	return response, metadata, nil
}
