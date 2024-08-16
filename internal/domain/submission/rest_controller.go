package submission

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/middleware"
	"github.com/highfive-compfest/seatudy-backend/internal/response"
)

type Controller struct {
	useCase *UseCase
}

func NewRestController(r *gin.Engine, uc *UseCase) {
	c := &Controller{useCase: uc}

	submissionGroup := r.Group("/v1/submissions")
	{
		submissionGroup.POST("", middleware.Authenticate(), middleware.RequireRole("student"), c.createSubmission)
		submissionGroup.GET("/:id", c.getSubmissionByID)
		submissionGroup.PUT("/:id", middleware.Authenticate(), middleware.RequireRole("student"), c.updateSubmission)
		submissionGroup.DELETE("/:id", middleware.Authenticate(), middleware.RequireRole("student"), c.deleteSubmission)
		submissionGroup.GET("/assignments/:assignmentId", c.getAllSubmissionsByAssignment)
		submissionGroup.PUT("/grade/:id", middleware.Authenticate(), middleware.RequireRole("instructor"), c.gradeSubmission)
	}

}

func (c *Controller) createSubmission(ctx *gin.Context) {
	var req CreateSubmissionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid submission data: "+err.Error(), nil).Send(ctx)
		return
	}

	userID, exists := ctx.Get("user.id")
	if !exists {
		response.NewRestResponse(http.StatusInternalServerError, "User ID not found in request context", nil).Send(ctx)
		return
	}

	if err := c.useCase.CreateSubmission(ctx, &req, userID.(string)); err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
		return
	}
	response.NewRestResponse(http.StatusCreated, "Submission created successfully", nil).Send(ctx)
}

func (c *Controller) getSubmissionByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid ID", nil).Send(ctx)
		return
	}

	submission, err := c.useCase.GetSubmissionByID(ctx, id)
	if err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
		return
	}
	response.NewRestResponse(http.StatusOK, "Submission retrieved successfully", submission).Send(ctx)
}

func (c *Controller) gradeSubmission(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid Submission ID", nil).Send(ctx)
		return
	}

	userID, exists := ctx.Get("user.id")
	if !exists {
		response.NewRestResponse(http.StatusInternalServerError, "User ID not found in request context", nil).Send(ctx)
		return
	}

	var req GradeSubmissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid grading data: "+err.Error(), nil).Send(ctx)
		return
	}

	if err := c.useCase.GradeSubmission(ctx, userID.(string), id, req.Grade); err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), "Failed to grade submission: "+err.Error(), apierror.GetPayload(err)).Send(ctx)
		return
	}

	response.NewRestResponse(http.StatusOK, "Submission graded successfully", nil).Send(ctx)
}

func (c *Controller) updateSubmission(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid ID", nil).Send(ctx)
		return
	}

	var req UpdateSubmissionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid submission data: "+err.Error(), nil).Send(ctx)
		return
	}

	userID, exists := ctx.Get("user.id")
	if !exists {
		response.NewRestResponse(http.StatusInternalServerError, "User ID not found in request context", nil).Send(ctx)
		return
	}

	if err := c.useCase.UpdateSubmission(ctx, id, &req, userID.(string)); err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
		return
	}
	response.NewRestResponse(http.StatusOK, "Submission updated successfully", nil).Send(ctx)
}

func (c *Controller) deleteSubmission(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid ID", nil).Send(ctx)
		return
	}

	userID, exists := ctx.Get("user.id")
	if !exists {
		response.NewRestResponse(http.StatusInternalServerError, "User ID not found in request context", nil).Send(ctx)
		return
	}

	if err := c.useCase.DeleteSubmission(ctx, id, userID.(string)); err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
		return
	}
	response.NewRestResponse(http.StatusOK, "Submission deleted successfully", nil).Send(ctx)
}

func (c *Controller) getAllSubmissionsByAssignment(ctx *gin.Context) {
	assignmentId, err := uuid.Parse(ctx.Param("assignmentId"))
	if err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid Assignment ID", nil).Send(ctx)
		return
	}

	submissions, err := c.useCase.GetAllSubmissionsByAssignment(ctx, assignmentId)
	if err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
		return
	}
	response.NewRestResponse(http.StatusOK, "Submissions retrieved successfully", submissions).Send(ctx)
}
