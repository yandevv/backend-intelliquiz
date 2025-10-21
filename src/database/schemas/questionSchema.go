package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Question struct {
	ID        string          `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	Content   string          `json:"content,omitempty" gorm:"not null"`
	QuizID    string          `json:"quiz_id,omitempty" gorm:"not null"`
	Quiz      *Quiz           `json:"quiz,omitempty"`
	Choices   []Choice        `json:"choices,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) (err error) {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return
}
