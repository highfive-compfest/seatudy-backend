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
	Price        float32          `json:"price" gorm:"type:numeric(11,2);not null;check:price >= 0"`
	ImageURL     string           `json:"image_url" gorm:"type:text"`
	SyllabusURL  string           `json:"syllabus_url" gorm:"type:text"`
	InstructorID uuid.UUID        `json:"instructor_id" gorm:"not null"`
	Difficulty   CourseDifficulty `json:"difficulty" gorm:"type:course_difficulty;not null"`
	CreatedAt    time.Time        `json:"created_at" gorm:"default:now();not null"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `json:"-" gorm:"index"`
}

type Role string

const (
	Student    Role = "student"
	Instructor Role = "instructor"
)

type User struct {
	ID              uuid.UUID      `json:"id" gorm:"primaryKey"`
	Email           string         `json:"email" gorm:"type:varchar(320);unique;not null"`
	IsEmailVerified bool           `json:"is_email_verified" gorm:"not null;default:false"`
	Name            string         `json:"name" gorm:"type:varchar(100);not null"`
	PasswordHash    string         `json:"-" gorm:"type:char(60);not null"`
	Role            Role           `json:"role" gorm:"type:user_role;not null"`
	ImageURL        string         `json:"image_url" gorm:"type:text"`
	Balance         float32        `json:"balance" gorm:"type:numeric(11,2);not null;default:0"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}