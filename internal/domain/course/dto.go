package course

import (
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"mime/multipart"
)

type CreateCourseRequest struct {
	Title       string                  `form:"title" binding:"required"`
	Description string                  `form:"description"`
	Price       float32                 `form:"price" binding:"required,gte=0"`
	Image       *multipart.FileHeader   `form:"image"`
	Syllabus    *multipart.FileHeader   `form:"syllabus"`
	Difficulty  schema.CourseDifficulty `form:"difficulty" binding:"required,oneof=beginner intermediate advanced expert"`
}

type UpdateCourseRequest struct {
	Title       *string                  `form:"title,omitempty"`
	Description *string                  `form:"description,omitempty"`
	Price       *float32                 `form:"price,omitempty" binding:"omitempty,gte=0"`
	Image       *multipart.FileHeader    `form:"image,omitempty"`    // Handled separately, not through direct JSON binding
	Syllabus    *multipart.FileHeader    `form:"syllabus,omitempty"` // Handled separately
	Difficulty  *schema.CourseDifficulty `form:"difficulty,omitempty" binding:"omitempty,oneof=beginner intermediate advanced expert"`
}
