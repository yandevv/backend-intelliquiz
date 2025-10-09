package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizScore struct {
	ID        string          `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	QuizID    string          `json:"quiz_id,omitempty" gorm:"not null"`
	Quiz      *Quiz           `json:"quiz,omitempty"`
	UserID    string          `json:"user_id,omitempty" gorm:"not null"`
	User      *User           `json:"user,omitempty"`
	Score     uint            `json:"score,omitempty" gorm:"not null"`
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
