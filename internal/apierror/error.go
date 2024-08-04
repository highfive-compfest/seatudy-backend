package apierror

import "net/http"

var (
	ErrInternalServer = ApiError{HttpStatus: http.StatusInternalServerError, Message: "INTERNAL_SERVER_ERROR"}
)
