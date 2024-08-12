package assignment

import (
	"net/http"

	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
)

var (
    ErrAssignmentNotFound        = apierror.ApiError{HttpStatus: http.StatusNotFound, Message: "ASSIGNMENT_NOT_FOUND"}
    ErrInvalidCourseData     = apierror.ApiError{HttpStatus: http.StatusBadRequest, Message: "INVALID_COURSE_DATA"}
    ErrNotOwnerAccess        = apierror.ApiError{HttpStatus: http.StatusUnauthorized, Message: "NOT_YOUR_COURSE"}
    ErrUnauthorizedAccess    = apierror.ApiError{HttpStatus: http.StatusUnauthorized, Message: "UNAUTHORIZED_ACCESS"}
    ErrForbiddenOperation    = apierror.ApiError{HttpStatus: http.StatusForbidden, Message: "FORBIDDEN_OPERATION"}
    ErrDatabaseOperationFail = apierror.ApiError{HttpStatus: http.StatusInternalServerError, Message: "DATABASE_OPERATION_FAILED"}
    ErrS3UploadFail          = apierror.ApiError{HttpStatus: http.StatusInternalServerError, Message: "S3_UPLOAD_FAILED"}
    ErrUUIDGenerationFail    = apierror.ApiError{HttpStatus: http.StatusInternalServerError, Message: "UUID_GENERATION_FAILED"}
    ErrEditConflict          = apierror.ApiError{HttpStatus: http.StatusConflict, Message: "EDIT_CONFLICT"}
)