package schema

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
	Price        int64          `json:"price" gorm:"not null"`
	Rating       float32          `json:"rating" gorm:"type:numeric(2,1);default:0.0;not null;check:rating >= 0.0 AND rating <= 5.0;index"`
	ReviewCount  int64            `json:"review_count" gorm:"type:bigint;default:0;not null"`
	ImageURL     string           `json:"image_url" gorm:"type:text"`
	SyllabusURL  string           `json:"syllabus_url" gorm:"type:text"`
	InstructorID uuid.UUID        `json:"instructor_id" gorm:"not null"`
	Difficulty   CourseDifficulty `json:"difficulty" gorm:"type:course_difficulty;not null"`
	Materials    []Material       `json:"materials" gorm:"foreignKey:CourseID"`
	CreatedAt    time.Time        `json:"created_at" gorm:"default:now();not null"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `json:"-" gorm:"index"`
}
