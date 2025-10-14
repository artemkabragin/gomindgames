package service

import (
	"mindgames/internal/domain"
)

type IUserService interface {
	Create(user *domain.User, password string) error
	GetByUsername(username string) (*domain.User, error)
}

type ITokenService interface {
	GenerateAccessToken(user domain.User) (string, error)
	GenerateRefreshToken(user domain.User) (domain.RefreshToken, error)
	GetRefreshByValue(refreshTokenValue string) (*domain.RefreshToken, error)
}
