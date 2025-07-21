package router

import (
	"github.com/MukizuL/vk-test/internal/controller"
	mw "github.com/MukizuL/vk-test/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func newRouter(c *controller.Controller, mw *mw.MiddlewareService) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	gin.Logger()

	api := router.Group("/api/v1")

	api.POST("/users", c.Register)
	api.POST("/tokens/authentication", c.Login)

	api.Use(mw.ShouldAuthorization()).GET("/ads", c.ListAds)
	api.Use(mw.MustAuthorization()).POST("/ads", c.CreateAd)

	return router
}

func Provide() fx.Option {
	return fx.Provide(newRouter)
}
