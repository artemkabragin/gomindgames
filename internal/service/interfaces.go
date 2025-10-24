package service

import (
	"context"
	"mindgames/internal/domain"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User, password string) error
	GetByUsername(username string) (*domain.User, error)
}

type TokenService interface {
	GenerateAccessToken(user domain.User) (string, error)
	GenerateRefreshToken(user domain.User) (domain.RefreshToken, error)
	GetRefreshByValue(refreshTokenValue string) (*domain.RefreshToken, error)
}
