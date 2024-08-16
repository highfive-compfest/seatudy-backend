package submission

import (
	"net/http"

	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
)

var (
	ErrAssignmentNotFound = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusNotFound).
				WithMessage("ASSIGNMENT_NOT_FOUND").
				Build()

	ErrInvalidCourseData = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusBadRequest).
				WithMessage("INVALID_COURSE_DATA").
				Build()

	ErrNotOwnerCourse = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusBadRequest).
				WithMessage("NOT_OWNER_ACCESS").
				Build()

	ErrNotEnrollCourse = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusBadRequest).
				WithMessage("NOT_ENROLL_ACCESS").
				Build()

	ErrSubmissionAlreadyExists = apierror.NewApiErrorBuilder().
					WithHttpStatus(http.StatusBadRequest).
					WithMessage("SUBMISSION_ALREADY_EXISTS_ACCESS").
					Build()

	ErrNotOwnerSubmission = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusBadRequest).
				WithMessage("NOT_OWNER_ACCESS").
				Build()

	ErrUnauthorizedAccess = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusUnauthorized).
				WithMessage("UNAUTHORIZED_ACCESS").
				Build()

	ErrForbiddenOperation = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusForbidden).
				WithMessage("FORBIDDEN_OPERATION").
				Build()

	ErrDatabaseOperationFail = apierror.NewApiErrorBuilder().
					WithHttpStatus(http.StatusInternalServerError).
					WithMessage("DATABASE_OPERATION_FAILED").
					Build()

	ErrS3UploadFail = apierror.NewApiErrorBuilder().
			WithHttpStatus(http.StatusInternalServerError).
			WithMessage("S3_UPLOAD_FAILED").
			Build()

	ErrUUIDGenerationFail = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusInternalServerError).
				WithMessage("UUID_GENERATION_FAILED").
				Build()

	ErrEditConflict = apierror.NewApiErrorBuilder().
			WithHttpStatus(http.StatusConflict).
			WithMessage("EDIT_CONFLICT").
			Build()
)
