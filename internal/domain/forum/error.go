package forum

import (
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"net/http"
)

var (
	ErrDiscussionNotFound = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusNotFound).
				WithMessage("DISCUSSION_NOT_FOUND")

	ErrReplyNotFound = apierror.NewApiErrorBuilder().
				WithHttpStatus(http.StatusNotFound).
				WithMessage("REPLY_NOT_FOUND")
)
