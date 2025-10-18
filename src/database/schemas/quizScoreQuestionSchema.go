package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizScoreQuestion struct {
	ID          string     `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	QuizScoreID string     `json:"quiz_score_id,omitempty" gorm:"not null"`
	QuizScore   *QuizScore `json:"quiz_score,omitempty"`
	QuestionID  string     `json:"question_id,omitempty" gorm:"not null"`
	Question    *Question  `json:"question,omitempty"`
	IsCorrect   bool       `json:"is_correct" gorm:"not null"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func (q *QuizScoreQuestion) BeforeCreate(tx *gorm.DB) (err error) {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return
}
