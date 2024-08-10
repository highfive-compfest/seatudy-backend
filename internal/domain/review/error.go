package review

import (
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"net/http"
)

var (
	ErrCourseAlreadyReviewed = apierror.ApiError{HttpStatus: http.StatusConflict, Message: "COURSE_ALREADY_REVIEWED"}
	ErrReviewNotFound        = apierror.ApiError{HttpStatus: http.StatusNotFound, Message: "REVIEW_NOT_FOUND"}
)
