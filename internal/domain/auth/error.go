package auth

import (
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"net/http"
)

var (
	ErrEmailAlreadyRegistered = apierror.NewApiErrorBuilder().
					WithHttpStatus(http.StatusConflict).
					WithMessage("EMAIL_ALREADY_REGISTERED")

	ErrInvalidCredentials = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusUnauthorized).
				WithMessage("INVALID_CREDENTIALS")

	ErrInvalidOTP = apierror.NewApiErrorBuilder().
			WithHttpStatus(http.StatusUnauthorized).
			WithMessage("INVALID_OTP")

	ErrExpiredOTP = apierror.NewApiErrorBuilder().
			WithHttpStatus(http.StatusUnauthorized).
			WithMessage("EXPIRED_OTP")

	ErrEmailAlreadyVerified = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusForbidden).
				WithMessage("EMAIL_ALREADY_VERIFIED")

	ErrInvalidResetPasswordLink = apierror.NewApiErrorBuilder().
					WithHttpStatus(http.StatusUnauthorized).
					WithMessage("INVALID_RESET_PASSWORD_LINK")

	ErrExpiredResetPasswordLink = apierror.NewApiErrorBuilder().
					WithHttpStatus(http.StatusUnauthorized).
					WithMessage("EXPIRED_RESET_PASSWORD_LINK")
)
