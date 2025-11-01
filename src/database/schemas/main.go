package schemas

import (
	"fmt"

	"gorm.io/gorm"
)

func Run(db *gorm.DB, fresh *bool) error {
	if *fresh {
		err := db.Migrator().DropTable(
			&User{},
			&Quiz{},
			&QuizUserLike{},
			&Question{},
			&Category{},
			&Game{},
			&GameQuestion{},
			&Choice{},
		)
		if err != nil {
			fmt.Println("Error dropping tables:", err)
			return err
		}
	}

	err := db.AutoMigrate(
		&User{},
		&Quiz{},
		&QuizUserLike{},
		&Question{},
		&Category{},
		&Game{},
		&GameQuestion{},
		&Choice{},
	)
	if err != nil {
		fmt.Println("Error during auto migration:", err)
		return err
	}

	err = db.SetupJoinTable(&Quiz{}, "UserLikes", &QuizUserLike{})
	if err != nil {
		fmt.Println("Error setting up join table:", err)
		return err
	}

	return nil
}
