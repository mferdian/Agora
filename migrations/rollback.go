package migrations

import (
	"Agora/model"

	"gorm.io/gorm"
)

func Rollback(db *gorm.DB) error {
	tables := []interface{}{
		&model.User{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	return nil
}
