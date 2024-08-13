package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Submission struct {
	ID           uuid.UUID      `json:"id" gorm:"primarykey"`
	AssignmentID uuid.UUID      `json:"assignment_id" gorm:"not null"`
	UserID       uuid.UUID      `json:"user_id" gorm:"not null"`
	Content      string         `json:"content" gorm:"type:varchar(1000)"`
	Grade        float64        `json:"grade" gorm:"type:numeric(4,1);check:grade BETWEEN 0 AND 100"`
	Attachments []Attachment   `json:"attachments" gorm:"foreignKey:SubmissionID"`
	CreatedAt    time.Time      `json:"created_at" gorm:"default:now();not null"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"" gorm:"index"`
}
