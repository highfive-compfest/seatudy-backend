package course

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseDifficulty string

const (
	Beginner     CourseDifficulty = "beginner"
	Intermediate CourseDifficulty = "intermediate"
	Advanced     CourseDifficulty = "advanced"
	Expert       CourseDifficulty = "expert"
)

type Course struct {
	ID           uuid.UUID        `json:"id" gorm:"primaryKey"`
	Title        string           `json:"title" gorm:"type:varchar(100);not null"`
	Description  string           `json:"description" gorm:"type:varchar(1000)"`
	Price        float32          `json:"price" gorm:"type:numeric(11,2);not null"`
	ImageURL     string           `json:"image_url" gorm:"type:text"`
	SyllabusURL  string           `json:"syllabus_url" gorm:"type:text"`
	InstructorID uuid.UUID        `json:"instructor_id" gorm:"not null"`
	Difficulty   CourseDifficulty `json:"difficulty" gorm:"type:course_difficulty;not null"`
	CreatedAt    time.Time        `json:"created_at" gorm:"default:now();not null"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `json:"-" gorm:"index"`
}
