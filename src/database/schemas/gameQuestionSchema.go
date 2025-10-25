package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameQuestion struct {
	ID         string     `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	GameID     string     `json:"game_id,omitempty" gorm:"not null"`
	Game       *Game      `json:"game,omitempty"`
	QuestionID string     `json:"question_id,omitempty" gorm:"not null"`
	Question   *Question  `json:"question,omitempty"`
	ChoiceID   *string    `json:"choice_id,omitempty"`
	Choice     *Choice    `json:"choice,omitempty"`
	Position   uint8      `json:"position" gorm:"not null"`
	AnsweredAt *time.Time `json:"answered_at,omitempty"`
	IsCorrect  bool       `json:"is_correct" gorm:"not null default:false"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

func (q *GameQuestion) BeforeCreate(tx *gorm.DB) (err error) {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return
}
