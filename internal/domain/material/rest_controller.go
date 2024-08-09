// file: material/rest_controller.go

package material

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"

	"github.com/highfive-compfest/seatudy-backend/internal/domain/course"
	"github.com/highfive-compfest/seatudy-backend/internal/middleware"
	"github.com/highfive-compfest/seatudy-backend/internal/response"
)

type RestController struct {
	useCase *UseCase
	courseUseCase *course.UseCase
	
}

func NewRestController(r *gin.Engine, uc *UseCase, cuc *course.UseCase) {
	c := &RestController{useCase: uc, courseUseCase: cuc}

	materialGroup := r.Group("/v1/materials")
    {
        materialGroup.POST("/",middleware.Authenticate(), c.create)
		materialGroup.GET("/:id", c.getByID)
		materialGroup.GET("/course/:id",c.getMaterialByCourse)
		materialGroup.GET("", c.getAll)
		materialGroup.PUT("/:id",middleware.Authenticate(), c.update)
		materialGroup.DELETE("/:id",middleware.Authenticate(), c.delete)
	}
	
}

func (c *RestController) create(ctx *gin.Context) {

	userRole , exists := ctx.Get("user.role")

	
	if !exists || userRole != "instructor" {
		response.NewRestResponse(http.StatusForbidden, "Only instructors are allowed to create material", nil).Send(ctx)
		return
	}

	var req CreateMaterialRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid material data: "+err.Error(), nil).Send(ctx)
		return
	}

	userID, exists := ctx.Get("user.id")
    if !exists {
        response.NewRestResponse(http.StatusInternalServerError, "User ID not found in request context", nil).Send(ctx)
        return
    }

    // Check if the course exists and is owned by the current user
	courseId, err := uuid.Parse(req.CourseID)

	if (err != nil) {
		err2 := apierror.ErrInternalServer
		response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
	}
	
	attachments, err := extractAttachments(ctx)

	if err != nil {
		response.NewRestResponse(http.StatusInternalServerError, "Failed to fetch attachment: "+err.Error(), nil).Send(ctx)
        return
	}

    req.Attachments = attachments
    course, err := c.courseUseCase.GetByID(ctx, courseId)
    if err != nil {
        response.NewRestResponse(http.StatusInternalServerError, "Failed to fetch material: "+err.Error(), nil).Send(ctx)
        return
    }

    if course.InstructorID.String() != userID {
        response.NewRestResponse(http.StatusForbidden, "Only the owner of the course can add materials", nil).Send(ctx)
        return
    }


	if err := c.useCase.CreateMaterial(ctx, req); err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
        return
	}
	response.NewRestResponse(http.StatusCreated, "Create material successfully", nil).Send(ctx)
}

func (c *RestController) getByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err = apierror.ErrInvalidParamId
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), err.Error()).Send(ctx)
		return
	}
	mat, err := c.useCase.GetMaterialByID(ctx, id)
	if err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
        return

	}
	response.NewRestResponse(http.StatusOK, "All Material Retrieve", mat).Send(ctx)
}

func (c *RestController) getAll(ctx *gin.Context) {
	mats, err := c.useCase.GetAllMaterials(ctx)
	if err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
        return
	}
	response.NewRestResponse(http.StatusOK, "All Material Retrieve", mats).Send(ctx)
}

func (c *RestController) getMaterialByCourse(ctx *gin.Context){
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err = apierror.ErrInvalidParamId
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), err.Error()).Send(ctx)
		return
	}

	course, err := c.courseUseCase.GetByID(ctx, id)
	if err != nil {
        response.NewRestResponse(http.StatusInternalServerError, "Failed to fetch material: "+err.Error(), nil).Send(ctx)
        return
    }

	response.NewRestResponse(http.StatusOK, "All Course Material Retrieve", course.Materials).Send(ctx)
	


}
func (c *RestController) update(ctx *gin.Context) {

	userRole , exists := ctx.Get("user.role")

	
	if !exists || userRole != "instructor" {
		response.NewRestResponse(http.StatusForbidden, "Only instructors are allowed to create material", nil).Send(ctx)
		return
	}

	var req UpdateMaterialRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewRestResponse(http.StatusBadRequest, "Invalid material data: "+err.Error(), nil).Send(ctx)
		return

	}

	attachments, err := extractAttachments(ctx)

	if err != nil {
		response.NewRestResponse(http.StatusInternalServerError, "Failed to fetch attachment: "+err.Error(), nil).Send(ctx)
        return
	}

    req.Attachments = attachments

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err = apierror.ErrInvalidParamId
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), err.Error()).Send(ctx)
		return
	}

	if err := c.useCase.UpdateMaterial(ctx, req, id); err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
            return
	}
	response.NewRestResponse(http.StatusCreated, "Update Material successfully", nil).Send(ctx)
}

func (c *RestController) delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err = apierror.ErrInvalidParamId
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), err.Error()).Send(ctx)
		return
	}
	if err := c.useCase.DeleteMaterial(ctx, id); err != nil {
		response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetDetail(err)).Send(ctx)
        return
	}
	response.NewRestResponse(http.StatusOK, "Delete Material successfully", nil).Send(ctx)
}

func extractAttachments(ctx *gin.Context) ([]AttachmentInput, error) {
    attachments := []AttachmentInput{}

    form, err := ctx.MultipartForm()
    if err != nil {
        return nil, fmt.Errorf("error retrieving multipart form: %w", err)
    }

    
    fileHeaders := form.File
    for i := 0; ; i++ {
        fileKey := fmt.Sprintf("attachments[%d][file]", i)
        descriptionKey := fmt.Sprintf("attachments[%d][description]", i)

        files, found := fileHeaders[fileKey]
        if !found {
            break 
        }

        description := ""
        if desc, ok := form.Value[descriptionKey]; ok && len(desc) > 0 {
            description = desc[0] 
        }

        if len(files) > 0 {
            attachments = append(attachments, AttachmentInput{
                File:        files[0], 
                Description: description,
            })
        }
    }

    return attachments, nil
}