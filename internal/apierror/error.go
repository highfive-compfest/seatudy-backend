package apierror

import "net/http"

var (
	ErrInternalServer   = ApiError{HttpStatus: http.StatusInternalServerError, Message: "INTERNAL_SERVER_ERROR"}
	ErrValidation       = ApiError{HttpStatus: http.StatusBadRequest, Message: "VALIDATION_ERROR"}
	ErrTokenEmpty       = ApiError{HttpStatus: http.StatusUnauthorized, Message: "TOKEN_EMPTY"}
	ErrTokenInvalid     = ApiError{HttpStatus: http.StatusUnauthorized, Message: "TOKEN_INVALID"}
	ErrTokenExpired     = ApiError{HttpStatus: http.StatusUnauthorized, Message: "TOKEN_EXPIRED"}
	ErrEmailNotVerified = ApiError{HttpStatus: http.StatusUnauthorized, Message: "EMAIL_NOT_VERIFIED"}
	ErrForbidden        = ApiError{HttpStatus: http.StatusForbidden, Message: "FORBIDDEN"}
)
