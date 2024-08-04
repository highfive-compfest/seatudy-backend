package response

import "github.com/gin-gonic/gin"

type Response struct {
	HttpCode int    `json:"-"`
	Message  string `json:"message"`
	Payload  any    `json:"payload"`
}

func NewRestResponse(httpCode int, message string, payload any) *Response {
	return &Response{
		HttpCode: httpCode,
		Message:  message,
		Payload:  payload,
	}
}

func (r *Response) Send(ctx *gin.Context) {
	ctx.JSON(r.HttpCode, r)
}
