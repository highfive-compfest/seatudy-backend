package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Assignment struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey"`
	CourseID    uuid.UUID      `json:"course_id" gorm:"not null"`
	Title       string         `json:"title" gorm:"type:varchar(150);not null"`
	Description string         `json:"description" gorm:"type:varchar(2000)"`
	Due         *time.Time     `json:"due"`
	Attachments []Attachment   `json:"attachments" gorm:"foreignKey:AssignmentID"`
	CreatedAt   time.Time      `json:"created_at" gorm:"default:now();not null"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"" gorm:"index"`
}