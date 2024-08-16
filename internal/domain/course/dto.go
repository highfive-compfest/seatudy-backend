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
	Category    schema.CourseCategory   `form:"category" binding:"required,oneof='Web Development' 'Game Development' 'Cloud Computing' 'Data Science & Analytics' 'Programming Languages' 'Cybersecurity' 'Mobile App Development' 'Database Management' 'Software Development' 'DevOps & Automation' 'Networking' 'AI & Machine Learning' 'Internet of Things (IoT)' 'Blockchain & Cryptocurrency' 'Augmented Reality (AR) & Virtual Reality (VR)'"`
}

type UpdateCourseRequest struct {
	Title       *string                  `form:"title,omitempty"`
	Description *string                  `form:"description,omitempty"`
	Price       *int64                 `form:"price,omitempty" binding:"omitempty,gte=0"`
	Image       *multipart.FileHeader    `form:"image,omitempty"`    // Handled separately, not through direct JSON binding
	Syllabus    *multipart.FileHeader    `form:"syllabus,omitempty"` // Handled separately
	Difficulty  *schema.CourseDifficulty `form:"difficulty,omitempty" binding:"omitempty,oneof=beginner intermediate advanced expert"`
	Category    *schema.CourseCategory   `form:"category" binding:"required,oneof='Web Development' 'Game Development' 'Cloud Computing' 'Data Science & Analytics' 'Programming Languages' 'Cybersecurity' 'Mobile App Development' 'Database Management' 'Software Development' 'DevOps & Automation' 'Networking' 'AI & Machine Learning' 'Internet of Things (IoT)' 'Blockchain & Cryptocurrency' 'Augmented Reality (AR) & Virtual Reality (VR)'"`
}

type CoursesPaginatedResponse struct {
    Courses    []schema.Course       `json:"courses"`
    Pagination pagination.Pagination `json:"pagination"`
}

type PaginationRequest struct {
    Page     int    `form:"page" binding:"required,min=1"`
	Limit    int    `form:"limit" binding:"required,min=1,max=30"`
}

type CourseProgress struct {
	UserId string `form:"user_id" binding:"required"`
}

type CourseProgressResponse struct {
	Course   schema.Course `json:"course"`
	Progress float64       `json:"progress"`
}

type SearchPaginationRequest struct {
	Title  string  `form:"title" binding:"required"` 
    Page   int     `form:"page" binding:"required,min=1"`
	Limit  int     `form:"limit" binding:"required,min=1,max=30"`
}

type FilterCoursesRequest struct {
    Rating    *float32              `form:"rating" binding:"omitempty,min=0,max=5"`
    Category  *string `form:"category" binding:"omitempty,oneof='Web Development' 'Game Development' 'Cloud Computing' 'Data Science & Analytics' 'Programming Languages' 'Cybersecurity' 'Mobile App Development' 'Database Management' 'Software Development' 'DevOps & Automation' 'Networking' 'AI & Machine Learning' 'Internet of Things (IoT)' 'Blockchain & Cryptocurrency' 'Augmented Reality (AR) & Virtual Reality (VR)'"`
    Difficulty *string `form:"difficulty" binding:"omitempty,oneof=beginner intermediate advanced expert"`
    Sort      *string               `form:"sort" binding:"omitempty,oneof=highest lowest"` 
    Page      int                  `form:"page" binding:"required,min=1"`
    Limit     int                  `form:"limit" binding:"required,min=1,max=50"`
}