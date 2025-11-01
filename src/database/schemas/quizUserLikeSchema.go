package schemas

import (
	"time"

	"gorm.io/gorm"
)

type QuizUserLike struct {
	QuizID    string    `json:"quiz_id" gorm:"primaryKey;not null"`
	UserID    string    `json:"user_id" gorm:"primaryKey;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}

func (q *QuizUserLike) BeforeCreate(tx *gorm.DB) (err error) {
	if q.CreatedAt.IsZero() {
		q.CreatedAt = time.Now()
	}

	return nil
}
