package assignment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/course"
	"github.com/highfive-compfest/seatudy-backend/internal/middleware"
	"github.com/highfive-compfest/seatudy-backend/internal/response"
)

type RestController struct {
	useCase       *UseCase
	courseUseCase *course.UseCase
}

func NewRestController(r *gin.Engine, uc *UseCase, cuc *course.UseCase) {
	c := &RestController{useCase: uc, courseUseCase: cuc}

	assignmentGroup := r.Group("/v1/assignments")
	{
		assignmentGroup.POST("", middleware.Authenticate(), middleware.RequireRole("instructor"), c.createAssignment)
		assignmentGroup.GET("/:id", c.getAssignmentByID)
		assignmentGroup.PUT("/:id", middleware.Authenticate(), middleware.RequireRole("instructor"), c.updateAssignment)
		assignmentGroup.POST("/addAttachment/:assignmentId", middleware.Authenticate(), middleware.RequireRole("instructor"), c.addAttachment)
		assignmentGroup.DELETE("/:id", middleware.Authenticate(), middleware.RequireRole("instructor"), c.deleteAssignment)
		assignmentGroup.GET("/course/:courseId", middleware.Authenticate(), c.getAssignmentsByCourse)
	}
}

func (c *RestController) createAssignment(ctx *gin.Context) {
	var req CreateAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid assignment data: "+err.Error(), nil).Send(ctx)
		return
	}

	userID, exists := ctx.Get("user.id")
	if !exists {
		response.NewRestResponse(http.StatusInternalServerError, "User ID not found in request context", nil).Send(ctx)
		return
	}

	courseID, err := uuid.Parse(req.CourseID)
	if err != nil {
		err2 := apierror.ErrInternalServer.Build()
		response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
	}

	course, err := c.courseUseCase.GetByID(ctx, courseID)
	if err != nil {
		response.NewRestResponse(http.StatusInternalServerError, "Failed to fetch course: "+err.Error(), nil).Send(ctx)
		return
	}

	if course.InstructorID.String() != userID {
		response.NewRestResponse(http.StatusForbidden, "Only the owner of the course can add assignment", nil).Send(ctx)
		return
	}

	if err := c.useCase.CreateAssignment(ctx, req, courseID); err != nil {
		response.NewRestResponse(http.StatusInternalServerError, "Failed to create assignment: "+err.Error(), nil).Send(ctx)
		return
	}
	response.NewRestResponse(http.StatusCreated, "Assignment created successfully", nil).Send(ctx)
}

func (c *RestController) addAttachment(ctx *gin.Context) {

	id, err := uuid.Parse(ctx.Param("assignmentId"))
	if err != nil {
		err = apierror.ErrInvalidParamId.Build()
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), err.Error()).Send(ctx)
		return
	}

	err = c.verifyAssignmentOwnership(ctx, id)
	if err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
		return
	}

	var req AttachmentInput
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid attachment data: "+err.Error(), nil).Send(ctx)
		return
	}

	if err := c.useCase.AddAttachment(ctx, id, req); err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
		return
	}
	response.NewRestResponse(http.StatusOK, "Add attachment successfully", nil).Send(ctx)
}

func (c *RestController) getAssignmentByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid ID", nil).Send(ctx)
		return
	}

	assignment, err := c.useCase.GetAssignmentByID(ctx, id)
	if err != nil {
		response.NewRestResponse(http.StatusInternalServerError, "Failed to fetch assignment: "+err.Error(), nil).Send(ctx)
		return
	}
	response.NewRestResponse(http.StatusOK, "Assignment retrieved successfully", assignment).Send(ctx)
}

func (c *RestController) updateAssignment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid ID", nil).Send(ctx)
		return
	}

	err = c.verifyAssignmentOwnership(ctx, id)
	if err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
		return
	}

	var req UpdateAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid assignment data: "+err.Error(), nil).Send(ctx)
		return
	}

	if err := c.useCase.UpdateAssignment(ctx, id, req); err != nil {
		response.NewRestResponse(http.StatusInternalServerError, "Failed to update assignment: "+err.Error(), nil).Send(ctx)
		return
	}
	response.NewRestResponse(http.StatusOK, "Assignment updated successfully", nil).Send(ctx)
}

func (c *RestController) deleteAssignment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid ID", nil).Send(ctx)
		return
	}

	err = c.verifyAssignmentOwnership(ctx, id)
	if err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
		return
	}

	if err := c.useCase.DeleteAssignment(ctx, id); err != nil {
		response.NewRestResponse(http.StatusInternalServerError, "Failed to delete assignment: "+err.Error(), nil).Send(ctx)
		return
	}
	response.NewRestResponse(http.StatusOK, "Assignment deleted successfully", nil).Send(ctx)
}

func (c *RestController) getAssignmentsByCourse(ctx *gin.Context) {
	courseId, err := uuid.Parse(ctx.Param("courseId"))
	if err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid Course ID", nil).Send(ctx)
		return
	}

	_, err = c.courseUseCase.GetByID(ctx, courseId)
	if err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
		return
	}

	assignments, err := c.useCase.GetAssignmentsByCourse(ctx, courseId)
	if err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
		return
	}
	response.NewRestResponse(http.StatusOK, "Assignments retrieved successfully", assignments).Send(ctx)
}

func (c *RestController) verifyAssignmentOwnership(ctx *gin.Context, assignmentId uuid.UUID) error {

	ass, err := c.useCase.GetAssignmentByID(ctx, assignmentId)
	if err != nil {
		return err
	}

	courseData, err := c.courseUseCase.GetByID(ctx, ass.CourseID)
	if err != nil {
		return err
	}
	instructorID, exists := ctx.Get("user.id")
	if !exists {
		return ErrUnauthorizedAccess.Build()
	}
	if courseData.InstructorID.String() != instructorID {
		return ErrNotOwnerAccess.Build()
	}

	return nil
}
