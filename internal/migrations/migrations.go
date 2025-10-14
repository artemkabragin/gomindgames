package migrations

import (
	"mindgames/internal/domain"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&domain.RefreshToken{}); err != nil {
		return err
	}

	return nil
}
