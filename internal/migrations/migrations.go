package migrations

import (
	"gorm.io/gorm"

	"mindgames/internal/model"
)

func Run(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.RefreshToken{}); err != nil {
		return err
	}

	return nil
}
