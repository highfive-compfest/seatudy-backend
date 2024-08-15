package review

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

	reviewGroup := engine.Group("/v1/reviews")
	{
		reviewGroup.POST("",
			middleware.Authenticate(),
			middleware.RequireRole("student"),
			controller.Create(),
		)
		reviewGroup.GET("",
			controller.Get(),
		)
		reviewGroup.PATCH("/:id",
			middleware.Authenticate(),
			middleware.RequireRole("student"),
			controller.Update(),
		)
		reviewGroup.DELETE("/:id",
			middleware.Authenticate(),
			middleware.RequireRole("student"),
			controller.Delete(),
		)
	}
}

func (c *RestController) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateReviewRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			err2 := apierror.ErrValidation.Build()
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		res, err := c.uc.Create(ctx, &req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusCreated, "CREATE_REVIEW_SUCCESS", res).Send(ctx)
	}
}

func (c *RestController) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req GetReviewsRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			err2 := apierror.ErrValidation.Build()
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		res, err := c.uc.Get(ctx, &req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "GET_REVIEWS_SUCCESS", res).Send(ctx)
	}
}

func (c *RestController) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req UpdateReviewRequest
		if err := ctx.ShouldBindUri(&req); err != nil {
			err2 := apierror.ErrValidation.Build()
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			err2 := apierror.ErrValidation.Build()
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		err := c.uc.Update(ctx, &req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "UPDATE_REVIEW_SUCCESS", nil).Send(ctx)
	}
}

func (c *RestController) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req DeleteReviewRequest
		if err := ctx.ShouldBindUri(&req); err != nil {
			err2 := apierror.ErrValidation.Build()
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		err := c.uc.Delete(ctx, &req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "DELETE_REVIEW_SUCCESS", nil).Send(ctx)
	}
}
