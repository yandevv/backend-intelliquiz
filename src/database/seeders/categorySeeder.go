package seeders

import (
	"intelliquiz/src/database/schemas"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CategoriesSeeding(db *gorm.DB) {
	categories := []string{
		"Ciências Gerais",
		"Matemática",
		"História",
		"Geografia",
		"Literatura",
		"Arte",
		"úsica",
		"Esportes",
		"Tecnologia",
		"Filmes e TV",
	}

	for _, name := range categories {
		db.Create(&schemas.Category{
			ID:   uuid.NewString(),
			Name: name,
		})
	}
}
