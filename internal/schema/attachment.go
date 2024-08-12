package schema

import (
	"github.com/google/uuid"
)

type Attachment struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey"`
	URL          string    `json:"url" gorm:"type:text;not null"`
	AssignmentID *uuid.UUID `json:"assignment_id" gorm:""`
	MaterialID   *uuid.UUID `json:"material_id" gorm:""`
	Description  string    `json:"description" gorm:"type:varchar(1000)"`
}
