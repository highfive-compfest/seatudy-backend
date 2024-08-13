package submission

import (
	"mime/multipart"


)

// CreateSubmissionRequest captures the necessary fields to create a submission
type CreateSubmissionRequest struct {
	AssignmentID string `form:"assignment_id" binding:"required" `
	Content      string    `form:"content" `
	// Attachments are not part of JSON, they are handled via multipart forms
	Attachments []*multipart.FileHeader `form:"attachments"`
}

// UpdateSubmissionRequest captures the fields that can be updated in a submission
type UpdateSubmissionRequest struct {
	Content      *string    `form:"content,omitempty"`
	// Attachments are not part of JSON, they are handled via multipart forms
	Attachments []*multipart.FileHeader `form:"attachments,omitempty"`
}

type GradeSubmissionRequest struct {
    Grade float64 `json:"grade" binding:"required,min=0,max=100"` // Ensure grade is between 0 and 100
}