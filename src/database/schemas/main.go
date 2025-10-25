package schemas

import "gorm.io/gorm"

func Run(db *gorm.DB, fresh *bool) error {
	if *fresh {
		db.Migrator().DropTable(
			&User{},
			&Quiz{},
			&Question{},
			&Category{},
			&Game{},
			&GameQuestion{},
			&Choice{},
		)
	}

	return db.AutoMigrate(
		&User{},
		&Quiz{},
		&Question{},
		&Category{},
		&Game{},
		&GameQuestion{},
		&Choice{},
	)
}
