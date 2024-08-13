package schema

import (
	"github.com/google/uuid"
	"time"
)

type ForumDiscussion struct {
	ID        uuid.UUID `json:"id" goorm:"primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null"`
	CourseID  uuid.UUID `json:"course_id" gorm:"not null,index"`
	Title     string    `json:"title" gorm:"type:varchar(150);not null"`
	Content   string    `json:"content" gorm:"type:varchar(30000);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-" gorm:"index"`
}

type ForumReply struct {
	ID                uuid.UUID `json:"id" gorm:"primaryKey"`
	UserID            uuid.UUID `json:"user_id" gorm:"not null"`
	ForumDiscussionID uuid.UUID `json:"forum_discussion_id" gorm:"not null;index"`
	CourseID          uuid.UUID `json:"course_id" gorm:"not null"`
	Content           string    `json:"content" gorm:"type:varchar(30000);not null"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	DeletedAt         time.Time `json:"-" gorm:"index"`
}
