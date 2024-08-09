package material

import (
	"time"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/attachment"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
type Material struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey"`
	CourseID    uuid.UUID      `json:"course_id" gorm:"not null"`
	Title       string         `json:"title" gorm:"type:varchar(150);not null"`
	Description string         `json:"description" gorm:"type:varchar(2000)"`
	Attachments []attachment.Attachment   `json:"attachments" gorm:"foreignKey:MaterialID"`
	CreatedAt   time.Time      `json:"created_at" gorm:"default:now();not null"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

}