package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Quiz struct {
	ID         string          `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	Name       string          `json:"name,omitempty" gorm:"size:60;not null"`
	CategoryID string          `json:"category_id,omitempty" gorm:"not null"`
	Category   *Category       `json:"category,omitempty"`
	CreatedBy  string          `json:"created_by,omitempty"`
	User       *User           `json:"user,omitempty" gorm:"foreignKey:CreatedBy"`
	Questions  []Question      `json:"questions,omitempty"`
	CreatedAt  *time.Time      `json:"created_at,omitempty"`
	UpdatedAt  *time.Time      `json:"updated_at,omitempty"`
	DeletedAt  *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (q *Quiz) BeforeCreate(tx *gorm.DB) (err error) {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return
}

func (q *Quiz) AfterDelete(tx *gorm.DB) (err error) {
	tx.Clauses(clause.Returning{}).Where("quiz_id = ?", q.ID).Delete(&Question{})
	return
}
