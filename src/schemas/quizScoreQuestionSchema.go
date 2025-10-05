package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizScoreQuestion struct {
	ID          string `json:"id" gorm:"type:uuid;primaryKey"`
	QuizScoreID string `json:"quiz_score_id" gorm:"not null"`
	QuizScore   QuizScore
	QuestionID  string `json:"question_id" gorm:"not null"`
	Question    Question
	IsCorrect   bool      `json:"is_correct" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (q *QuizScoreQuestion) BeforeCreate(tx *gorm.DB) (err error) {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return
}
