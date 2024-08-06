package user

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

	userGroup := engine.Group("/v1/users")
	{
		userGroup.GET("/me",
			middleware.Authenticate(),
			controller.GetMe(),
		)
		userGroup.PATCH("/me",
			middleware.Authenticate(),
			controller.Update(),
		)
		userGroup.GET("/:id",
			controller.GetByID(),
		)
	}
}

func (c *RestController) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req GetUserByIDRequest
		if err := ctx.ShouldBindUri(&req); err != nil {
			err2 := apierror.ErrValidation
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		res, err := c.uc.GetByID(&req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "GET_USER_SUCCESS", res).Send(ctx)
	}
}

func (c *RestController) GetMe() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := GetUserByIDRequest{ID: ctx.Value("user.id").(string)}

		res, err := c.uc.GetByID(&req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "GET_USER_SUCCESS", res).Send(ctx)
	}
}

func (c *RestController) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req UpdateUserRequest
		if err := ctx.ShouldBind(&req); err != nil {
			err2 := apierror.ErrValidation
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		if err := c.uc.Update(ctx, &req); err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "UPDATE_USER_SUCCESS", nil).Send(ctx)
	}
}
