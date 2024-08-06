package auth

import (
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"net/http"
)

var (
	ErrEmailAlreadyRegistered   = apierror.ApiError{HttpStatus: http.StatusConflict, Message: "EMAIL_ALREADY_REGISTERED"}
	ErrInvalidCredentials       = apierror.ApiError{HttpStatus: http.StatusUnauthorized, Message: "INVALID_CREDENTIALS"}
	ErrInvalidOTP               = apierror.ApiError{HttpStatus: http.StatusUnauthorized, Message: "INVALID_OTP"}
	ErrExpiredOTP               = apierror.ApiError{HttpStatus: http.StatusUnauthorized, Message: "EXPIRED_OTP"}
	ErrEmailAlreadyVerified     = apierror.ApiError{HttpStatus: http.StatusForbidden, Message: "EMAIL_ALREADY_VERIFIED"}
	ErrInvalidResetPasswordLink = apierror.ApiError{HttpStatus: http.StatusUnauthorized, Message: "INVALID_RESET_PASSWORD_LINK"}
	ErrExpiredResetPasswordLink = apierror.ApiError{HttpStatus: http.StatusUnauthorized, Message: "EXPIRED_RESET_PASSWORD_LINK"}
)
