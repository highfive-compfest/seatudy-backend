package course

import (
	"mime/multipart"

	"github.com/highfive-compfest/seatudy-backend/internal/pagination"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
)

type CreateCourseRequest struct {
	Title       string                  `form:"title" binding:"required"`
	Description string                  `form:"description"`
	Price       int64                 `form:"price" binding:"required,gte=0"`
	Image       *multipart.FileHeader   `form:"image"`
	Syllabus    *multipart.FileHeader   `form:"syllabus"`
	Difficulty  schema.CourseDifficulty `form:"difficulty" binding:"required,oneof=beginner intermediate advanced expert"`
}

type UpdateCourseRequest struct {
	Title       *string                  `form:"title,omitempty"`
	Description *string                  `form:"description,omitempty"`
	Price       *int64                 `form:"price,omitempty" binding:"omitempty,gte=0"`
	Image       *multipart.FileHeader    `form:"image,omitempty"`    // Handled separately, not through direct JSON binding
	Syllabus    *multipart.FileHeader    `form:"syllabus,omitempty"` // Handled separately
	Difficulty  *schema.CourseDifficulty `form:"difficulty,omitempty" binding:"omitempty,oneof=beginner intermediate advanced expert"`
}

type CoursesPaginatedResponse struct {
    Courses    []schema.Course       `json:"courses"`
    Pagination pagination.Pagination `json:"pagination"`
}

type PaginationRequest struct {
    Page     int    `form:"page" binding:"required,min=1"`
	Limit    int    `form:"limit" binding:"required,min=1,max=30"`
}