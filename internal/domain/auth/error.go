package auth

import (
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"net/http"
)

var (
	ErrEmailAlreadyRegistered = apierror.ApiError{HttpStatus: http.StatusConflict, Message: "EMAIL_ALREADY_REGISTERED"}
	ErrInvalidCredentials     = apierror.ApiError{HttpStatus: http.StatusUnauthorized, Message: "INVALID_CREDENTIALS"}
	ErrTokenEmpty             = apierror.ApiError{HttpStatus: http.StatusUnauthorized, Message: "TOKEN_EMPTY"}
	ErrTokenInvalid           = apierror.ApiError{HttpStatus: http.StatusUnauthorized, Message: "TOKEN_INVALID"}
	ErrTokenExpired           = apierror.ApiError{HttpStatus: http.StatusUnauthorized, Message: "TOKEN_EXPIRED"}
)
