package assignment

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type CreateAssignmentRequest struct {
	CourseID    string `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Due         *time.Time `json:"due,omitempty"`
}

type UpdateAssignmentRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Due         *time.Time `json:"due,omitempty"`
}

type AssignmentResponse struct {
	ID          uuid.UUID `json:"id"`
	CourseID    uuid.UUID `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Due         time.Time `json:"due"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}


type AttachmentInput struct {
    File        *multipart.FileHeader `form:"file" binding:"required"` // The actual file
    Description string                `form:"description"`             // Optional description
}