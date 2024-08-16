package user

import (
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"net/http"
)

var (
	ErrUserNotFound = apierror.NewApiErrorBuilder().
		WithHttpStatus(http.StatusNotFound).
		WithMessage("USER_NOT_FOUND")
)
