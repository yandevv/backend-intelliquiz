package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Game struct {
	ID        string          `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	UserID    string          `json:"user_id,omitempty" gorm:"not null"`
	User      *User           `json:"user,omitempty"`
	QuizID    string          `json:"quiz_id,omitempty" gorm:"not null"`
	Quiz      *Quiz           `json:"quiz,omitempty"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (g *Game) BeforeCreate(tx *gorm.DB) (err error) {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return
}
