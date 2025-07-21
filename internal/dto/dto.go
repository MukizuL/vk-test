package dto

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/MukizuL/vk-test/internal/validator"
	"github.com/greatcloak/decimal"
)

type AuthFormRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func ValidateAuthFormRequest(v *validator.Validator, req AuthFormRequest) {
	v.Check(utf8.RuneCountInString(req.Login) >= 3, "login", "must be minimum 3 symbols")
	v.Check(utf8.RuneCountInString(req.Login) <= 255, "login", "must be not greater than 255 symbols")

	v.Check(utf8.RuneCountInString(req.Password) >= 8, "password", "must be minimum 8 symbols")
	v.Check(len(req.Password) <= 72, "password", fmt.Sprintf("too long: %d bytes", len(req.Password)))
}

type GetAdsResponse struct {
	Login       string          `json:"login"`
	Owned       bool            `json:"owned,omitempty"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	ImageURL    string          `json:"image_url"`
	Price       decimal.Decimal `json:"price"`
	CreatedAt   time.Time       `json:"created_at"`
}

type CreateAdRequest struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	ImageURL    string          `json:"image_URL"`
	Price       decimal.Decimal `json:"price"`
}

func ValidateCreateAdRequest(v *validator.Validator, req CreateAdRequest) {
	v.Check(utf8.RuneCountInString(req.Title) >= 3, "title", "must be minimum 3 symbols")
	v.Check(utf8.RuneCountInString(req.Title) <= 255, "title", "must be not greater than 255 symbols")

	v.Check(utf8.RuneCountInString(req.Description) <= 10_000, "description", "must be not greater than 10 000 symbols")

	v.Check(req.Price.Exponent() >= -2, "price", "too precise")
}

type CreateAdResponse struct {
	Login       string          `json:"login"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	ImageURL    string          `json:"image_url"`
	Price       decimal.Decimal `json:"price"`
	CreatedAt   time.Time       `json:"created_at"`
}
