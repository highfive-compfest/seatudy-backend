package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	Student    Role = "student"
	Instructor Role = "instructor"
)

type User struct {
	ID              uuid.UUID      `json:"id" gorm:"primaryKey"`
	Email           string         `json:"email" gorm:"type:varchar(320);unique;not null"`
	IsEmailVerified bool           `json:"is_email_verified" gorm:"not null;default:false"`
	Name            string         `json:"name" gorm:"type:varchar(50);not null"`
	PasswordHash    string         `json:"-" gorm:"type:char(60);not null"`
	Role            Role           `json:"role" gorm:"type:user_role;not null"`
	ImageURL        string         `json:"image_url" gorm:"type:text"`
	Balance         float32        `json:"balance" gorm:"type:numeric(11,2);not null;default:0"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

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
	Materials    []Material             `json:"materials" gorm:"foreignKey:CourseID"`
	CreatedAt    time.Time        `json:"created_at" gorm:"default:now();not null"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `json:"-" gorm:"index"`
}

type Material struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey"`
	CourseID    uuid.UUID      `json:"course_id" gorm:"not null"`
	Title       string         `json:"title" gorm:"type:varchar(150);not null"`
	Description string         `json:"description" gorm:"type:varchar(2000)"`
	Attachments []Attachment   `json:"attachments" gorm:"foreignKey:MaterialID"`
	CreatedAt   time.Time      `json:"created_at" gorm:"default:now();not null"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

}

type Attachment struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey"`
	URL         string         `json:"url" gorm:"type:text;not null"`
	MaterialID  uuid.UUID      `json:"material_id" gorm:"not null"`
	Description string         `json:"description" gorm:"type:varchar(1000)"`
}