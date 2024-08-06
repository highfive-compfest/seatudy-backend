package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/middleware"
	"github.com/highfive-compfest/seatudy-backend/internal/response"
	"net/http"
)

type RestController struct {
	uc *UseCase
}

func NewRestController(engine *gin.Engine, uc *UseCase) {
	controller := &RestController{uc: uc}

	authGroup := engine.Group("/v1/auth")
	{
		authGroup.POST("/register", controller.Register())
		authGroup.POST("/login", controller.Login())
		authGroup.POST("/refresh", controller.Refresh())
		authGroup.POST("/verification/email/send",
			middleware.Authenticate(),
			controller.SendOTP(),
		)
		authGroup.PATCH("/verification/email/verify",
			middleware.Authenticate(),
			controller.VerifyOTP(),
		)
		authGroup.POST("/password/reset/request", controller.SendResetPasswordLink())
		authGroup.PATCH("/password/reset/verify", controller.ResetPassword())
		authGroup.PATCH("/password/change",
			middleware.Authenticate(),
			controller.ChangePassword(),
		)
	}

}

func (c *RestController) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req RegisterRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			err2 := apierror.ErrValidation
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		if err := c.uc.Register(&req); err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusCreated, "REGISTER_SUCCESS", nil).Send(ctx)
	}
}

func (c *RestController) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			err2 := apierror.ErrValidation
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		resp, err := c.uc.Login(&req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "LOGIN_SUCCESS", resp).Send(ctx)
	}
}

func (c *RestController) Refresh() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req RefreshRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			err2 := apierror.ErrValidation
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		resp, err := c.uc.Refresh(&req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "REFRESH_SUCCESS", resp).Send(ctx)
	}
}

func (c *RestController) SendOTP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := c.uc.SendOTP(ctx); err != nil {
			err2 := apierror.ErrValidation
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "OTP_SEND_SUCCESS", nil).Send(ctx)
	}
}

func (c *RestController) VerifyOTP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req VerifyEmailRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			err2 := apierror.ErrValidation
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		if err := c.uc.VerifyOTP(ctx, &req); err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "OTP_VERIFY_SUCCESS", nil).Send(ctx)
	}
}

func (c *RestController) SendResetPasswordLink() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req SendResetPasswordLinkRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			err2 := apierror.ErrValidation
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		if err := c.uc.SendResetPasswordLink(ctx, &req); err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "RESET_PASSWORD_LINK_SEND_SUCCESS", nil).Send(ctx)
	}
}

func (c *RestController) ResetPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ResetPasswordRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			err2 := apierror.ErrValidation
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		if err := c.uc.ResetPassword(ctx, &req); err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "RESET_PASSWORD_SUCCESS", nil).Send(ctx)
	}
}

func (c *RestController) ChangePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ChangePasswordRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			err2 := apierror.ErrValidation
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		if err := c.uc.ChangePassword(ctx, &req); err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "CHANGE_PASSWORD_SUCCESS", nil).Send(ctx)
	}
}
