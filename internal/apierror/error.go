package apierror

import "net/http"

var (
	ErrInternalServer = NewApiErrorBuilder().
				WithHttpStatus(http.StatusInternalServerError).
				WithMessage("INTERNAL_SERVER_ERROR")

	ErrValidation = NewApiErrorBuilder().
			WithHttpStatus(http.StatusBadRequest).
			WithMessage("VALIDATION_ERROR")

	ErrTokenEmpty = NewApiErrorBuilder().
			WithHttpStatus(http.StatusUnauthorized).
			WithMessage("TOKEN_EMPTY")

	ErrTokenInvalid = NewApiErrorBuilder().
			WithHttpStatus(http.StatusUnauthorized).
			WithMessage("TOKEN_INVALID")

	ErrTokenExpired = NewApiErrorBuilder().
			WithHttpStatus(http.StatusUnauthorized).
			WithMessage("TOKEN_EXPIRED")

	ErrEmailNotVerified = NewApiErrorBuilder().
				WithHttpStatus(http.StatusUnauthorized).
				WithMessage("EMAIL_NOT_VERIFIED")

	ErrForbidden = NewApiErrorBuilder().
			WithHttpStatus(http.StatusForbidden).
			WithMessage("FORBIDDEN")

	ErrNotYourResource = NewApiErrorBuilder().
				WithHttpStatus(http.StatusForbidden).
				WithMessage("NOT_YOUR_RESOURCE")

	ErrFileTooLarge = NewApiErrorBuilder().
			WithHttpStatus(http.StatusBadRequest).
			WithMessage("FILE_TOO_LARGE")

	ErrInvalidFileType = NewApiErrorBuilder().
				WithHttpStatus(http.StatusBadRequest).
				WithMessage("INVALID_FILE_TYPE")

	ErrInvalidParamId = NewApiErrorBuilder().
				WithHttpStatus(http.StatusBadRequest).
				WithMessage("INVALID_PARAM_ID")

	ErrInsufficientBalance = NewApiErrorBuilder().
				WithHttpStatus(http.StatusBadRequest).
				WithMessage("INSUFFICIENT_BALANCE")
)
