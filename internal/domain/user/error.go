package user

import (
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"net/http"
)

var (
	ErrUserNotFound = apierror.ApiError{HttpStatus: http.StatusNotFound, Message: "USER_NOT_FOUND"}
)
