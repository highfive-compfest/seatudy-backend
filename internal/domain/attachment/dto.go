package attachment

import "mime/multipart"

type AttachmentUpdateRequest struct {
	File        *multipart.FileHeader `form:"file" binding:"required"` // The actual file
	Description string                `form:"description"`             // Optional description
}