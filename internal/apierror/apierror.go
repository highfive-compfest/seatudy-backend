package apierror

import (
	"errors"
	"net/http"
)

type ApiError struct {
	HttpStatus int
	Message    string
	Payload    any
}

func (e ApiError) Error() string {
	return e.Message
}

func GetHttpStatus(err error) int {
	var httpErr ApiError
	if errors.As(err, &httpErr) {
		return httpErr.HttpStatus
	}
	return http.StatusInternalServerError
}

func AddPayload(apiErr ApiError, detail any) ApiError {
	apiErr.Payload = detail
	return apiErr
}

func GetDetail(err error) any {
	var httpErr ApiError
	if errors.As(err, &httpErr) {
		return httpErr.Payload
	}
	return nil
}
