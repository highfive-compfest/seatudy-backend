package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/response"
	"net/http"
)

type RestController struct {
	uc *UseCase
}

func NewRestController(router *gin.Engine, uc *UseCase) {
	controller := &RestController{uc: uc}

	authGroup := router.Group("/v1/auth")
	{
		authGroup.POST("/register", controller.Register())
		authGroup.POST("/login", controller.Login())
		authGroup.POST("/refresh", controller.Refresh())
		authGroup.POST("/verify-email", controller.VerifyEmail())
	}

}

func (c *RestController) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req RegisterRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.NewRestResponse(http.StatusBadRequest, "VALIDATION_ERROR", err.Error()).Send(ctx)
			return
		}

		err := c.uc.Register(req)
		if err != nil {
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
			response.NewRestResponse(http.StatusBadRequest, "VALIDATION_ERROR", err.Error()).Send(ctx)
			return
		}

		resp, err := c.uc.Login(req)
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
			response.NewRestResponse(http.StatusBadRequest, "VALIDATION_ERROR", err.Error()).Send(ctx)
			return
		}

		resp, err := c.uc.Refresh(req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "REFRESH_SUCCESS", resp).Send(ctx)
	}
}

func (c *RestController) VerifyEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Implementation here
	}
}
