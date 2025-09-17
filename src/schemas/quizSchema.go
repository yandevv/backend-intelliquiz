package schemas

import (
	"time"
)

type Quiz struct {
	ID         string    `json:"id" gorm:"default:uuid_generate_v3();primaryKey"`
	Name       string    `json:"name" gorm:"not null"`
	CategoryID string    `json:"category_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
