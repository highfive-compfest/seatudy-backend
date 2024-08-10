package review

import "github.com/google/uuid"

type CreateReviewRequest struct {
	CourseID uuid.UUID `json:"course_id" binding:"required,uuid"`
	Rating   int       `json:"rating" binding:"required,min=1,max=5"`
	Feedback string    `json:"feedback" binding:"required,max=255"`
}

type CreateReviewResponse struct {
	ID uuid.UUID `json:"id"`
}

type GetReviewsRequest struct {
	CourseID string `form:"course_id" binding:"required,uuid"`
	Rating   int    `form:"rating" binding:"omitempty,min=1,max=5"`
	Page     int    `form:"page" binding:"required,min=1"`
	Limit    int    `form:"limit" binding:"required,min=1,max=30"`
}

type UpdateReviewRequest struct {
	ID       string `uri:"id" binding:"required,uuid"`
	Rating   int    `json:"rating" binding:"omitempty,min=1,max=5"`
	Feedback string `json:"feedback" binding:"max=255"`
}

type DeleteReviewRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}
