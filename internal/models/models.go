package models

import (
	"time"

	"github.com/greatcloak/decimal"
)

type User struct {
	ID           string
	Login        string
	PasswordHash string
}

type Ad struct {
	ID          string
	UserLogin   string
	Title       string
	Description string
	ImageURL    string
	Price       decimal.Decimal
	CreatedAt   time.Time
}
