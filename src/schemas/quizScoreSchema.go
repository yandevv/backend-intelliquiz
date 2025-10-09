package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizScore struct {
	ID        string `json:"id" gorm:"type:uuid;primaryKey"`
	QuizID    string `json:"quiz_id" gorm:"not null"`
	Quiz      Quiz
	UserID    string `json:"user_id" gorm:"not null"`
	User      User
	Score     uint            `json:"score" gorm:"not null"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (q *QuizScore) BeforeCreate(tx *gorm.DB) (err error) {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return
}
