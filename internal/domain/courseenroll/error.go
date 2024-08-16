package courseenroll

import (
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"net/http"
)

var (
	ErrNotEnrolled = apierror.NewApiErrorBuilder().
		WithHttpStatus(http.StatusBadRequest).
		WithMessage("COURSE_NOT_ENROLLED")
)
