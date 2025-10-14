package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	ID                   uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Username             string         `json:"username"`
	PasswordHash         string         `json:"passwordHash"`
	IsOnboardingComplete bool           `json:"isOnboardingComplete"`
	RefreshTokens        []RefreshToken `json:"-" gorm:"foreignKey:UserID"`
}

type UserPublic struct {
	Username             string `json:"username"`
	IsOnboardingComplete bool   `json:"isOnboardingComplete"`
}

type UserClaims struct {
	UserID   string `json:"uid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
