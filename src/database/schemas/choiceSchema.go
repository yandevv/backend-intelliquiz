package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Choice struct {
	ID         string          `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	QuestionID string          `json:"question_id,omitempty" gorm:"not null"`
	Question   *Question       `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	Content    string          `json:"content,omitempty" gorm:"not null"`
	CreatedAt  *time.Time      `json:"created_at,omitempty"`
	UpdatedAt  *time.Time      `json:"updated_at,omitempty"`
	DeletedAt  *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (c *Choice) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}
