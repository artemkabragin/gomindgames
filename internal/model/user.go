package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                   uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username             string         `json:"username" gorm:"uniqueIndex;not null"`
	PasswordHash         string         `json:"passwordHash" gorm:"not null"`
	IsOnboardingComplete bool           `json:"isOnboardingComplete"`
	RefreshTokens        []RefreshToken `json:"-" gorm:"foreignKey:UserID"`
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"userId" gorm:"type:uuid;index;not null"`
	Value     string    `json:"value" gorm:"not null"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
}
