package courseenroll

import (
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"net/http"
)

var (
	ErrNotEnrolled = apierror.ApiError{HttpStatus: http.StatusForbidden, Message: "NOT_ENROLLED", Payload: nil}
)
