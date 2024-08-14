package schema

import (
	"github.com/google/uuid"
	"time"
)

type Notification struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	UserID    uuid.UUID `json:"-" gorm:"not null;index"`
	Title     string    `json:"title" gorm:"type:varchar(255);not null"`
	Detail    string    `json:"detail" gorm:"type:text"`
	IsRead    bool      `json:"is_read" gorm:"not null;default:false;index"`
	CreatedAt time.Time `json:"created_at"`
}
