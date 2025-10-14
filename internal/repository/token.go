package repository

import (
	"errors"
	"fmt"
	"mindgames/internal/domain"
	"time"

	"gorm.io/gorm"
)

type TokenRepositoryImpl struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) ITokenRepository {
	return &TokenRepositoryImpl{
		db: db,
	}
}

func (r TokenRepositoryImpl) Create(refreshToken *domain.RefreshToken) error {
	if err := r.db.Create(&refreshToken).Error; err != nil {
		return err
	}
	return nil
}

func (r TokenRepositoryImpl) GetRefreshByValue(refreshTokenValue string) (*domain.RefreshToken, error) {
	var refreshToken domain.RefreshToken
	if err := r.db.Preload("User").Where("value = ? AND expires_at > ?", refreshTokenValue, time.Now()).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &domain.RefreshToken{}, fmt.Errorf("invalid or expired refresh token: %w", err)
		}
		return &domain.RefreshToken{}, fmt.Errorf("database error: %w", err)
	}
	return &refreshToken, nil
}
