package schema


import (
    "time"

    "github.com/google/uuid"
)

type CourseEnroll struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
    UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
    CourseID  uuid.UUID `json:"course_id" gorm:"type:uuid;not null;index"`
    CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
}