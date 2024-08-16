package submission

import (
	"mime/multipart"


)


type CreateSubmissionRequest struct {
	AssignmentID string `form:"assignment_id" binding:"required" `
	Content      string    `form:"content" `
	Attachments []*multipart.FileHeader `form:"attachments"`
}


type UpdateSubmissionRequest struct {
	Content      *string    `form:"content,omitempty"`
	Attachments []*multipart.FileHeader `form:"attachments,omitempty"`
}

type GradeSubmissionRequest struct {
    Grade float64 `json:"grade" binding:"required,min=0,max=100"` 
}