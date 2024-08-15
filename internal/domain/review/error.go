package review

import (
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"net/http"
)

var (
	ErrCourseAlreadyReviewed = apierror.NewApiErrorBuilder().
					WithHttpStatus(http.StatusConflict).
					WithMessage("COURSE_ALREADY_REVIEWED")

	ErrReviewNotFound = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusNotFound).
				WithMessage("REVIEW_NOT_FOUND")
)
