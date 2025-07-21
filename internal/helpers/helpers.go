package helpers

import (
	"net/url"
	"strconv"

	"github.com/MukizuL/vk-test/internal/validator"
)

func ReadInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func ReadFloat64(qs url.Values, key string, defaultValue float64, v *validator.Validator) float64 {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		v.AddError(key, "must be a float value")
		return defaultValue
	}

	return i
}

func ReadString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}
