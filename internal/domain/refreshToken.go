package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `json:"userId" gorm:"type:uuid;index"`
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expiresAt"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
}

func (r *RefreshToken) BeforeCreate(*gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	if r.ExpiresAt.IsZero() {
		r.ExpiresAt = time.Now().Add(7 * 24 * time.Hour) // 7 days
	}
	return nil
}
