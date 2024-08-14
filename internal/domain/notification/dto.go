package notification

type GetByUserIDRequest struct {
	Limit int `form:"limit" binding:"required,max=30"`
	Page  int `form:"page" binding:"required"`
}
