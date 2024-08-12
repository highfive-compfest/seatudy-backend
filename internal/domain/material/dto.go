package material

import (
	"mime/multipart"


)

type AttachmentInput struct {
    File        *multipart.FileHeader `form:"file" binding:"required"` // The actual file
    Description string                `form:"description"`             // Optional description
}

type CreateMaterialRequest struct {
    CourseID    string          `form:"course_id" binding:"required"`
    Title       string             `form:"title" binding:"required"`
    Description string             `form:"description"`
    
}

type UpdateMaterialRequest struct {
    Title       *string             `form:"title"`
    Description *string             `form:"description"`
    
}