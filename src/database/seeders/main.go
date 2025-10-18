package seeders

import (
	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	CategoriesSeeding(db)
}
