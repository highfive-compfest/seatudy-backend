package schema

import (
	"github.com/google/uuid"
	"time"
)

type Review struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null;index:idx_user_course,unique"`
	CourseID  uuid.UUID `json:"course_id" gorm:"not null;index:idx_user_course,unique"`
	Rating    int       `json:"rating" gorm:"type:smallint;not null;check:rating >= 1 AND rating <= 5;index"`
	Feedback  string    `json:"feedback" gorm:"type:varchar(255); not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-" gorm:"index"`
}
