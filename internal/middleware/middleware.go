package middleware

import (
	"net/http"
	"strings"

	"github.com/MukizuL/vk-test/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type MiddlewareService struct {
	service *services.Services
}

func NewMiddlewareService(service *services.Services) *MiddlewareService {
	return &MiddlewareService{
		service: service,
	}
}

func Provide() fx.Option {
	return fx.Provide(NewMiddlewareService)
}

func (s *MiddlewareService) MustAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken := ctx.GetHeader("Authorization")
		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Error": "No access token",
			})
			return
		}

		userID, err := s.service.ValidateToken(strings.TrimPrefix(accessToken, "Bearer "))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Error": "Access token is invalid",
			})
			return
		}

		ctx.Set("userID", userID)

		ctx.Next()
	}
}

func (s *MiddlewareService) ShouldAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken := ctx.GetHeader("Authorization")
		if accessToken == "" {
			ctx.Set("userID", "")
			ctx.Next()
			return
		}

		userID, err := s.service.ValidateToken(strings.TrimPrefix(accessToken, "Bearer "))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Error": "Access token is invalid",
			})
			return
		}

		ctx.Set("userID", userID)

		ctx.Next()
	}
}
