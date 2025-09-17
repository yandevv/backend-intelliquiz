package schemas

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `json:"id" gorm:"default:uuid_generate_v3();primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
