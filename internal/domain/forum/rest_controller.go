package forum

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

	discussionsGroup := engine.Group("/v1/forums/discussions")
	discussionsGroup.Use(middleware.Authenticate())
	{
		discussionsGroup.POST("", controller.CreateDiscussion())
		discussionsGroup.GET("/:id", controller.GetDiscussionByID())
		discussionsGroup.GET("", controller.GetDiscussionsByCourseID())
		discussionsGroup.PATCH("/:id", controller.UpdateDiscussion())
		discussionsGroup.DELETE("/:id", controller.DeleteDiscussion())
	}

	repliesGroup := engine.Group("/v1/forums/replies")
	repliesGroup.Use(middleware.Authenticate())
	{
		repliesGroup.POST("", controller.CreateReply())
		repliesGroup.GET("/:id", controller.GetReplyByID())
		repliesGroup.GET("", controller.GetRepliesByDiscussionID())
		repliesGroup.PATCH("/:id", controller.UpdateReply())
		repliesGroup.DELETE("/:id", controller.DeleteReply())
	}
}

func (c *RestController) CreateDiscussion() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateForumDiscussionRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			err2 := apierror.ErrValidation.Build()
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		err := c.uc.CreateDiscussion(ctx, &req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusCreated, "CREATE_DISCUSSION_SUCCESS", nil).Send(ctx)
	}
}

func (c *RestController) GetDiscussionByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		discussion, err := c.uc.GetDiscussionByID(ctx, id)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "GET_DISCUSSION_SUCCESS", discussion).Send(ctx)
	}
}

func (c *RestController) GetDiscussionsByCourseID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req GetForumDiscussionsRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			err2 := apierror.ErrValidation.Build()
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		res, err := c.uc.GetDiscussionsByCourseID(ctx, &req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "GET_DISCUSSIONS_SUCCESS", res).Send(ctx)
	}
}

func (c *RestController) UpdateDiscussion() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req UpdateForumDiscussionRequest
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

		err := c.uc.UpdateDiscussion(ctx, &req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "UPDATE_DISCUSSION_SUCCESS", nil).Send(ctx)
	}
}

func (c *RestController) DeleteDiscussion() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := c.uc.DeleteDiscussion(ctx, id)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "DELETE_DISCUSSION_SUCCESS", nil).Send(ctx)
	}
}

func (c *RestController) CreateReply() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateForumReplyRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			err2 := apierror.ErrValidation.Build()
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		err := c.uc.CreateReply(ctx, &req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusCreated, "CREATE_REPLY_SUCCESS", nil).Send(ctx)
	}
}

func (c *RestController) GetReplyByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		reply, err := c.uc.GetReplyByID(ctx, id)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "GET_REPLY_SUCCESS", reply).Send(ctx)
	}
}

func (c *RestController) GetRepliesByDiscussionID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req GetForumRepliesRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			err2 := apierror.ErrValidation.Build()
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		res, err := c.uc.GetRepliesByDiscussionID(ctx, &req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "GET_REPLIES_SUCCESS", res).Send(ctx)
	}
}

func (c *RestController) UpdateReply() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req UpdateForumReplyRequest
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

		err := c.uc.UpdateReply(ctx, &req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "UPDATE_REPLY_SUCCESS", nil).Send(ctx)
	}
}

func (c *RestController) DeleteReply() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := c.uc.DeleteReply(ctx, id)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "DELETE_REPLY_SUCCESS", nil).Send(ctx)
	}
}
