package service

import (
	"fmt"
	"mindgames/internal/domain"
	"mindgames/internal/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenServiceImpl struct {
	r repository.TokenRepository
}

func NewTokenService(repo repository.TokenRepository) *TokenServiceImpl {
	return &TokenServiceImpl{
		repo,
	}
}

func (s TokenServiceImpl) GenerateAccessToken(user domain.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me-dev-secret"
	}

	now := time.Now()
	claims := domain.UserClaims{
		UserID:   user.ID.String(),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s TokenServiceImpl) GenerateRefreshToken(user domain.User) (domain.RefreshToken, error) {
	refresh := domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Value:     uuid.NewString(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := s.r.Create(&refresh); err != nil {
		return domain.RefreshToken{}, fmt.Errorf("error creating user: %w", err)
	}

	return refresh, nil
}

func (s TokenServiceImpl) GetRefreshByValue(refreshTokenValue string) (*domain.RefreshToken, error) {
	refresh, err := s.r.GetRefreshByValue(refreshTokenValue)

	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return refresh, nil
}
