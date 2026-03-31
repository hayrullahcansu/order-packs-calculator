package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TimeAwareEntity provides common fields for all database tables
type TimeAwareEntity struct {
	ID        uuid.UUID `gorm:"primaryKey;type:char(36);uniqueIndex" json:"id" example:"550e8400-e29b-41d4-a716-446655440000"` //is the primary key using uuid with char(36) storage and unique index.
	CreatedAt time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`                                                     // records when the entity was first created.
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-15T10:30:00Z"`                                                     // tracks the last modification time, automatically managed by gorm.
}

func (t *TimeAwareEntity) BeforeCreate(tx *gorm.DB) (err error) {
	uid := uuid.New()
	t.ID = uid
	return
}
