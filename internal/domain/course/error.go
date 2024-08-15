package course

import (
	"net/http"

	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
)

var (
	ErrCourseNotFound = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusNotFound).
				WithMessage("COURSE_NOT_FOUND")

	ErrInvalidCourseData = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusBadRequest).
				WithMessage("INVALID_COURSE_DATA")

	ErrUnauthorizedAccess = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusUnauthorized).
				WithMessage("UNAUTHORIZED_ACCESS")

	ErrNotOwnerAccess = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusUnauthorized).
				WithMessage("NOT_YOUR_COURSE")

	ErrForbiddenOperation = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusForbidden).
				WithMessage("FORBIDDEN_OPERATION")

	ErrDatabaseOperationFail = apierror.NewApiErrorBuilder().
					WithHttpStatus(http.StatusInternalServerError).
					WithMessage("DATABASE_OPERATION_FAILED")

	ErrS3UploadFail = apierror.NewApiErrorBuilder().
			WithHttpStatus(http.StatusInternalServerError).
			WithMessage("S3_UPLOAD_FAILED")

	ErrUUIDGenerationFail = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusInternalServerError).
				WithMessage("UUID_GENERATION_FAILED")

	ErrEditConflict = apierror.NewApiErrorBuilder().
			WithHttpStatus(http.StatusConflict).
			WithMessage("EDIT_CONFLICT")

	ErrAlreadyEnrolled = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusBadRequest).
				WithMessage("COURSE_ALREADY_ENROLLED")
)
