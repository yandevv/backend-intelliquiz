package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Quiz struct {
	ID         string    `json:"id" gorm:"type:uuid;primaryKey"`
	Name       string    `json:"name" gorm:"not null"`
	CategoryID string    `json:"category_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (q *Quiz) BeforeCreate(tx *gorm.DB) (err error) {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return
}
