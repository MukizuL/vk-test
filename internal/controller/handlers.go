package controller

import (
	"context"
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/MukizuL/vk-test/internal/dto"
	"github.com/MukizuL/vk-test/internal/errs"
	"github.com/MukizuL/vk-test/internal/filters"
	"github.com/MukizuL/vk-test/internal/helpers"
	"github.com/MukizuL/vk-test/internal/validator"
	"github.com/gin-gonic/gin"
)

func (c *Controller) Register(ctx *gin.Context) {
	ctxTO, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var data dto.AuthFormRequest
	err := ctx.BindJSON(&data)
	if err != nil {
		errs.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	v := validator.New()

	if dto.ValidateAuthFormRequest(v, data); !v.Valid() {
		errs.FailedValidationResponse(ctx, v.Errors)
		return
	}

	userID, err := c.service.CreateUser(ctxTO, data.Login, data.Password)
	if err != nil {
		if errors.Is(err, errs.ErrConflictLogin) {
			errs.ErrorResponse(ctx, http.StatusConflict, err.Error())
			return
		}

		errs.ServerErrorResponse(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user_id": userID,
	})
}

func (c *Controller) Login(ctx *gin.Context) {
	ctxTO, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var data dto.AuthFormRequest
	err := ctx.BindJSON(&data)
	if err != nil {
		errs.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	v := validator.New()

	if dto.ValidateAuthFormRequest(v, data); !v.Valid() {
		errs.FailedValidationResponse(ctx, v.Errors)
		return
	}

	token, err := c.service.LoginUser(ctxTO, data.Login, data.Password)
	if err != nil {
		if errors.Is(err, errs.ErrNotAuthorized) {
			errs.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		if errors.Is(err, errs.ErrUserNotFound) {
			errs.NotFoundResponse(ctx)
			return
		}

		errs.ServerErrorResponse(ctx, err.Error())
		return
	}

	ctx.Header("Authorization", "Bearer "+token)

	ctx.JSON(http.StatusOK, gin.H{
		"message": http.StatusText(http.StatusOK),
	})
}

func (c *Controller) ListAds(ctx *gin.Context) {
	ctxTO, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var input filters.Filters

	v := validator.New()

	qs := ctx.Request.URL.Query()

	input.Page = helpers.ReadInt(qs, "page", 1, v)
	input.PageSize = helpers.ReadInt(qs, "page_size", 20, v)

	input.Min = helpers.ReadFloat64(qs, "min", 0, v)
	input.Max = helpers.ReadFloat64(qs, "max", math.MaxInt32, v)

	input.Sort = helpers.ReadString(qs, "sort", "date")

	input.SortSafelist = []string{"created_at", "-created_at", "date", "-date", "price", "-price"}

	if filters.ValidateFilters(v, input); !v.Valid() {
		errs.FailedValidationResponse(ctx, v.Errors)
		return
	}

	value, exists := ctx.Get("userID")
	if !exists {
		panic("userID does not exist") // It must exist
	}

	userID, ok := value.(string)
	if !ok {
		panic("userID is not a string") // It means I made a mistake
	}

	ads, metadata, err := c.service.GetAds(ctxTO, userID, input)
	if err != nil {
		switch {
		default:
			errs.ServerErrorResponse(ctx, err.Error())
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ads": ads, "metadata": metadata})
}

func (c *Controller) CreateAd(ctx *gin.Context) {
	ctxTO, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var data dto.CreateAdRequest
	err := ctx.BindJSON(&data)
	if err != nil {
		errs.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	v := validator.New()

	if dto.ValidateCreateAdRequest(v, data); !v.Valid() {
		errs.FailedValidationResponse(ctx, v.Errors)
		return
	}

	value, exists := ctx.Get("userID")
	if !exists {
		panic("userID does not exist") // It must exist
	}

	userID, ok := value.(string)
	if !ok {
		panic("userID is not a string") // It means I made a mistake
	}

	ad, err := c.service.CreateAd(ctxTO, userID, &data)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			errs.NotFoundResponse(ctx)
		default:
			errs.ServerErrorResponse(ctx, err.Error())
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"ad": ad})
}
