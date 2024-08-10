package course

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
	uc *UseCase
}

func NewRestController(router *gin.Engine, uc *UseCase) {

	controller := &RestController{uc: uc}

	courseGroup := router.Group("/v1/courses")
	{
		courseGroup.GET("", controller.GetAll())
		courseGroup.GET("/:id", controller.GetByID())
		courseGroup.POST("", middleware.Authenticate(), controller.Create())
		courseGroup.PUT("/:id", middleware.Authenticate(), controller.Update())
		courseGroup.GET("/instructor/:id", middleware.Authenticate(), controller.GetInstructorCourse())
		courseGroup.DELETE("/:id",
			middleware.Authenticate(),
			middleware.RequireRole("instructor"),
			controller.Delete(),
		)
	}

}

func (c *RestController) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		courses, err := c.uc.GetAll(ctx)
		if err != nil {
			response.NewRestResponse(http.StatusInternalServerError, "Failed to retrieve courses", nil).Send(ctx)
			return
		}
		response.NewRestResponse(http.StatusOK, "Courses retrieved successfully", courses).Send(ctx)
	}
}

func (c *RestController) GetInstructorCourse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		instructorID, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			response.NewRestResponse(http.StatusBadRequest, "Invalid Instructor ID", nil).Send(ctx)
			return
		}

		// Fetch all courses by the instructor ID
		courses, err := c.uc.GetByInstructorID(ctx, instructorID)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
			return
		}

		if len(courses) == 0 {
			response.NewRestResponse(http.StatusOK, "No courses found for this instructor", nil).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "Courses retrieved successfully", courses).Send(ctx)
	}
}

func (c *RestController) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			response.NewRestResponse(http.StatusBadRequest, "Invalid ID", nil).Send(ctx)
			return
		}

		course, err := c.uc.GetByID(ctx, id)
		if err != nil {
			response.NewRestResponse(http.StatusInternalServerError, err.Error(), nil).Send(ctx)
			return
		}
		response.NewRestResponse(http.StatusOK, "Course retrieved successfully", course).Send(ctx)
	}
}

func (c *RestController) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userRole, exists := ctx.Get("user.role")

		log.Print(exists)
		if !exists || userRole != "instructor" {

			response.NewRestResponse(http.StatusForbidden, "Only instructors are allowed to create courses", nil).Send(ctx)
			return
		}

		var req CreateCourseRequest
		if err := ctx.ShouldBind(&req); err != nil {
			response.NewRestResponse(http.StatusBadRequest, "Invalid course data: "+err.Error(), nil).Send(ctx)
			return
		}

		imageFile, errImage := ctx.FormFile("image")
		if errImage != nil && errImage != http.ErrMissingFile {
			response.NewRestResponse(http.StatusBadRequest, "Could not retrieve image file", nil).Send(ctx)
			return
		}

		syllabusFile, errSyllabus := ctx.FormFile("syllabus")
		if errSyllabus != nil && errSyllabus != http.ErrMissingFile {
			response.NewRestResponse(http.StatusBadRequest, "Could not retrieve syllabus file", nil).Send(ctx)
			return
		}

		instructorID, exists := ctx.Get("user.id")
		if !exists {
			response.NewRestResponse(http.StatusInternalServerError, "Failed to get instructor ID from context", nil).Send(ctx)
			return
		}

		err := c.uc.Create(ctx, req, imageFile, syllabusFile, instructorID.(string))
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusCreated, "Course created successfully", nil).Send(ctx)
	}
}

func (c *RestController) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Ensure user role is "instructor"
		userRole, exists := ctx.Get("user.role")

		log.Print(exists)
		if !exists || userRole != "instructor" {
			// Logging for debugging
			log.Printf("Access denied or role not found: %v", userRole)
			response.NewRestResponse(http.StatusForbidden, "Only instructors are allowed to create courses", nil).Send(ctx)
			return
		}

		// Parse UUID from the URL parameter
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			response.NewRestResponse(http.StatusBadRequest, "Invalid ID", nil).Send(ctx)
			return
		}

		// Bind JSON payload to UpdateCourseRequest struct
		var req UpdateCourseRequest
		if err := ctx.ShouldBind(&req); err != nil {
			response.NewRestResponse(http.StatusBadRequest, "Invalid course data: "+err.Error(), nil).Send(ctx)
			return
		}

		// Handle optional file uploads
		imageFile, _ := ctx.FormFile("image")
		syllabusFile, _ := ctx.FormFile("syllabus")

		// Call the use case to update the course
		updatedCourse, err := c.uc.Update(ctx.Request.Context(), req, id, imageFile, syllabusFile)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
			return
		}

		// Successfully updated the course
		response.NewRestResponse(http.StatusOK, "Course updated successfully", updatedCourse).Send(ctx)
	}
}

func (c *RestController) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			response.NewRestResponse(http.StatusBadRequest, "Invalid ID", nil).Send(ctx)
			return
		}

		err = c.uc.Delete(ctx, id)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
			return
		}
		response.NewRestResponse(http.StatusOK, "Course deleted successfully", nil).Send(ctx)
	}
}
