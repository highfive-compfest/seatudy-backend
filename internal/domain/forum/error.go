package forum

import (
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"net/http"
)

var (
	ErrDiscussionNotFound = apierror.ApiError{HttpStatus: http.StatusNotFound, Message: "DISCUSSION_NOT_FOUND"}
	ErrReplyNotFound      = apierror.ApiError{HttpStatus: http.StatusNotFound, Message: "REPLY_NOT_FOUND"}
)
