package attachment

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/middleware"
	"github.com/highfive-compfest/seatudy-backend/internal/response"
)

type RestController struct {
	useCase *UseCase
}

func NewRestController(r *gin.Engine, uc *UseCase) {
	c := &RestController{useCase: uc}

	attachmentsGroup := r.Group("/v1/attachments")
	{
		attachmentsGroup.PUT("/:id", middleware.Authenticate(), c.update)
		attachmentsGroup.GET("/:id", middleware.Authenticate(), c.get)
		attachmentsGroup.DELETE("/:id", middleware.Authenticate(), c.delete)

	}
}

func (c *RestController) get(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err = apierror.ErrInvalidParamId.Build()
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), err.Error()).Send(ctx)
		return
	}
	attachment, err := c.useCase.GetAttachmentByID(ctx, id)
	if err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
		return
	}
	response.NewRestResponse(http.StatusOK, "Retrieve attachment successfully", attachment).Send(ctx)
}

func (c *RestController) update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err = apierror.ErrInvalidParamId.Build()
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), err.Error()).Send(ctx)
		return
	}

	var req AttachmentUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid attachment data: "+err.Error(), nil).Send(ctx)
		return

	}
	log.Println(req.File)

	attachment, err := c.useCase.UpdateAttachment(ctx, id, req)
	if err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
		return
	}
	response.NewRestResponse(http.StatusCreated, "Update attachment successfully", attachment).Send(ctx)
}

func (c *RestController) delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err = apierror.ErrInvalidParamId.Build()
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), err.Error()).Send(ctx)
		return
	}
	if err := c.useCase.DeleteAttachment(ctx, id); err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
		return
	}
	response.NewRestResponse(http.StatusCreated, "Delete attachment successfully", nil).Send(ctx)
}
