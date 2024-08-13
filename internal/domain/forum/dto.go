package forum

import "github.com/google/uuid"

type CreateForumDiscussionRequest struct {
	CourseID uuid.UUID `json:"course_id" binding:"required,uuid"`
	Title    string    `json:"title" binding:"required,max=150"`
	Content  string    `json:"content" binding:"required,max=30000"`
}

type GetForumDiscussionsRequest struct {
	CourseID string `form:"course_id" binding:"required,uuid"`
	Page     int    `form:"page" binding:"required"`
	Limit    int    `form:"limit" binding:"required,max=30"`
}

type UpdateForumDiscussionRequest struct {
	ID      string `uri:"id" binding:"required,uuid"`
	Title   string `json:"title" binding:"max=150"`
	Content string `json:"content" binding:"max=30000"`
}

type DeleteForumDiscussionRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type CreateForumReplyRequest struct {
	DiscussionID uuid.UUID `json:"discussion_id" binding:"required,uuid"`
	Content      string    `json:"content" binding:"required,max=30000"`
}

type GetForumRepliesRequest struct {
	DiscussionID string `uri:"discussion_id" binding:"required,uuid"`
	Page         int    `form:"page" binding:"required"`
	Limit        int    `form:"limit" binding:"required,max=30"`
}

type UpdateForumReplyRequest struct {
	ID      string `uri:"id" binding:"required,uuid"`
	Content string `json:"content" binding:"max=30000"`
}
