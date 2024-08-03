package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	Student    Role = "student"
	Instructor Role = "instructor"
)

type User struct {
	ID              uuid.UUID `json:"id" gorm:"primaryKey"`
	Email           string    `json:"email" gorm:"type:varchar(320);unique;not null"`
	IsEmailVerified bool      `json:"is_email_verified" gorm:"type:boolean;not null;default:false"`
	Name            string    `json:"name" gorm:"type:varchar(100);not null"`
	PasswordHash    string    `json:"-" gorm:"type:char(60);not null"`
	Role            Role      `json:"role" gorm:"type:user_role;not null"`
	ImageURL        string    `json:"image_url" gorm:"type:text"`
	Balance         int       `json:"balance" gorm:"type:numeric(11,2);not null;default:0"`
	gorm.Model
}
