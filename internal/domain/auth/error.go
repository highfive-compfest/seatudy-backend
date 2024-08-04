package auth

import (
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"net/http"
)

var (
	ErrEmailAlreadyRegistered = apierror.ApiError{HttpStatus: http.StatusConflict, Message: "EMAIL_ALREADY_REGISTERED"}
	ErrInvalidCredentials     = apierror.ApiError{HttpStatus: http.StatusUnauthorized, Message: "INVALID_CREDENTIALS"}
	ErrInvalidToken           = apierror.ApiError{HttpStatus: http.StatusUnauthorized, Message: "INVALID_TOKEN"}
)
