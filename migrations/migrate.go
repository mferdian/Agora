package migrations

import (
	"Agora/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Proposal{},
	); err != nil {
		return err
	}

	return nil
}
