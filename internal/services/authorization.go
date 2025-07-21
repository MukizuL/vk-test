package services

import (
	"context"

	"github.com/MukizuL/vk-test/internal/errs"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser Creates a user with given login and password. Returns a userID and an error.
func (s *Services) CreateUser(ctx context.Context, login, password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", errs.ErrInternalServerError
	}

	userID, err := s.storage.CreateNewUser(ctx, login, string(passwordHash))
	if err != nil {
		s.logger.Error("failed to create a user", zap.String("login", login), zap.Error(err))
		return "", errs.TransformPGErrors(err)
	}

	return userID, nil
}

// LoginUser Logs in a user with given login and password. Returns a JWT and an error.
func (s *Services) LoginUser(ctx context.Context, login, password string) (string, error) {
	user, err := s.storage.GetUserByLogin(ctx, login)
	if err != nil {
		s.logger.Error("failed to get a user by login", zap.String("login", login), zap.Error(err))
		return "", errs.TransformPGErrors(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errs.ErrNotAuthorized
	}

	accessTokenSigned, err := s.CreateToken(user.ID)
	if err != nil {
		return "", err
	}

	return accessTokenSigned, nil
}
