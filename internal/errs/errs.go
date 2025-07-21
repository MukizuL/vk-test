package errs

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrConflictLogin           = errors.New("this login is already used by other user")
	ErrUserNotFound            = errors.New("login or password is incorrect")
	ErrInternalServerError     = errors.New("internal server error")
	ErrNotAuthorized           = errors.New("invalid token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)

func ErrorResponse(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{"error": message})
}

func ServerErrorResponse(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": message})
}

func FailedValidationResponse(ctx *gin.Context, errors map[string]string) {
	ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": errors})
}

func NotFoundResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "resource not found"})
}

func TransformPGErrors(err error) error {
	var pgErr *pgconn.PgError

	switch {
	case errors.As(err, &pgErr):
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return ErrConflictLogin
		}
	case errors.Is(err, pgx.ErrNoRows):
		return ErrUserNotFound
	default:
		return ErrInternalServerError
	}

	return ErrInternalServerError
}
