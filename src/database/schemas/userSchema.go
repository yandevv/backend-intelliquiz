package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string          `json:"id,omitempty" gorm:"type:uuid;primaryKey"`
	Username  string          `json:"username,omitempty" gorm:"size:50;uniqueIndex;not null"`
	Password  string          `json:"password,omitempty" gorm:"size:50;not null"`
	Email     string          `json:"email,omitempty" gorm:"size:254;uniqueIndex;not null"`
	Name      string          `json:"name,omitempty" gorm:"size:60;not null"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
