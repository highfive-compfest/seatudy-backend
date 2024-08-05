package apierror

import "net/http"

var (
	ErrInternalServer = ApiError{HttpStatus: http.StatusInternalServerError, Message: "INTERNAL_SERVER_ERROR"}
	ErrTokenEmpty     = ApiError{HttpStatus: http.StatusUnauthorized, Message: "TOKEN_EMPTY"}
	ErrTokenInvalid   = ApiError{HttpStatus: http.StatusUnauthorized, Message: "TOKEN_INVALID"}
	ErrTokenExpired   = ApiError{HttpStatus: http.StatusUnauthorized, Message: "TOKEN_EXPIRED"}
)
