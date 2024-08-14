package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/middleware"
	"github.com/highfive-compfest/seatudy-backend/internal/response"
	"net/http"
)

type RestController struct {
	useCase *UseCase
}

func NewRestController(engine *gin.Engine, useCase *UseCase) {
	controller := &RestController{useCase: useCase}

	notificationsGroup := engine.Group("/v1/notifications")
	notificationsGroup.Use(middleware.Authenticate())
	{
		notificationsGroup.GET("/my", controller.GetMy())
		notificationsGroup.GET("/my/unread-count", controller.GetUnreadCount())
		notificationsGroup.PATCH("/read/:id", controller.UpdateRead())
	}
}

func (c *RestController) GetMy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req GetByUserIDRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			err2 := apierror.ErrValidation
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		notifications, err := c.useCase.GetMy(ctx, &req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "GET_NOTIFICATIONS_SUCCESS", notifications).Send(ctx)
	}
}

func (c *RestController) GetUnreadCount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		count, err := c.useCase.GetUnreadCount(ctx)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "GET_UNREAD_COUNT_SUCCESS", count).Send(ctx)
	}
}

func (c *RestController) UpdateRead() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			err := apierror.ErrInvalidParamId
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
			return
		}

		err := c.useCase.UpdateRead(id)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "UPDATE_READ_SUCCESS", nil).Send(ctx)
	}
}
