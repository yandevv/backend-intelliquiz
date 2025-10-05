package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Quiz struct {
	ID         string     `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	Name       string     `json:"name,omitempty" gorm:"not null"`
	CategoryID string     `json:"category_id,omitempty" gorm:"not null"`
	Category   *Category  `json:"category,omitempty"`
	CreatedBy  string     `json:"created_by,omitempty"`
	User       *User      `json:"user,omitempty" gorm:"foreignKey:CreatedBy"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

func (q *Quiz) BeforeCreate(tx *gorm.DB) (err error) {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return
}
