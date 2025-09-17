package schemas

import (
	"time"

	"gorm.io/gorm"
)

type QuizScore struct {
	ID        string         `json:"id" gorm:"default:uuid_generate_v3();primaryKey"`
	QuizID    string         `json:"quiz_id" gorm:"not null"`
	UserID    string         `json:"user_id" gorm:"not null"`
	Score     uint           `json:"score" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
