package schemas

import (
	"time"
)

type QuizScoreQuestion struct {
	ID          string    `json:"id" gorm:"default:uuid_generate_v3();primaryKey"`
	QuizScoreID string    `json:"quiz_score_id" gorm:"not null"`
	QuestionID  string    `json:"question_id" gorm:"not null"`
	IsCorrect   bool      `json:"is_correct" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
