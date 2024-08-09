package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Role string

const (
	Student    Role = "student"
	Instructor Role = "instructor"
)

type User struct {
	ID              uuid.UUID      `json:"id" gorm:"primaryKey"`
	Email           string         `json:"email" gorm:"type:varchar(320);unique;not null;index:,type:hash"`
	IsEmailVerified bool           `json:"is_email_verified" gorm:"not null;default:false"`
	Name            string         `json:"name" gorm:"type:varchar(50);not null"`
	PasswordHash    string         `json:"-" gorm:"type:char(60);not null"`
	Role            Role           `json:"role" gorm:"type:user_role;not null"`
	ImageURL        string         `json:"image_url" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}
