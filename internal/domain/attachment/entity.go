package attachment

import (
		"github.com/google/uuid"
)


type Attachment struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey"`
	URL         string         `json:"url" gorm:"type:text;not null"`
	MaterialID  uuid.UUID      `json:"material_id" gorm:"not null"`
	Description string         `json:"description" gorm:"type:varchar(1000)"`
}