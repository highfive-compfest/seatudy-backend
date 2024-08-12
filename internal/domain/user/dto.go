package user

import (
	"mime/multipart"
)

type GetUserByIDRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type GetUserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	Role     string `json:"role"`
}

type UpdateUserRequest struct {
	Name      string                `form:"name" binding:"max=50"`
	ImageFile *multipart.FileHeader `form:"image_file"`
}
