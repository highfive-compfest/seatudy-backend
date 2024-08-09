package course

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"slices"

	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/config"
	"github.com/highfive-compfest/seatudy-backend/internal/fileutil"
)

type UseCase struct {
    courseRepo Repository
}

func NewUseCase(courseRepo Repository) *UseCase {
    return &UseCase{courseRepo: courseRepo}
}

func (uc *UseCase) GetAll(ctx context.Context) ([]Course, error) {
    return uc.courseRepo.GetAll(ctx)
}

func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (Course, error) {
    return uc.courseRepo.GetByID(ctx, id)
}

func (uc *UseCase) Create(ctx context.Context, req CreateCourseRequest, imageFile, syllabusFile *multipart.FileHeader, instructorID string) error {
    var imageUrl, syllabusUrl string
    var err error

	uuidInstructorID, err := uuid.Parse(instructorID)
    if err != nil {
        log.Printf("Error parsing instructor ID: %v", err)
        return ErrUnauthorizedAccess // Or any other appropriate error
    }

    
	id, err := uuid.NewV7()
	if err != nil {
		log.Println("Error generating UUID: ", err)
		return apierror.ErrInternalServer
	}

    // Upload image if present
    if imageFile != nil {

        if imageFile.Size > 2*fileutil.MegaByte {
			err2 := apierror.ErrFileTooLarge
			apierror.AddPayload(&err2, map[string]string{
				"max_size":      "2 MB",
				"received_size": fileutil.ByteToAppropriateUnit(imageFile.Size),
			})
			return err2
		}
        fileType, err := fileutil.DetectMultipartFileType(imageFile)

        if err != nil {
			log.Println("Error detecting image type: ", err)
			return apierror.ErrInternalServer
		}

        allowedTypes := fileutil.ImageContentTypes
		if !slices.Contains(allowedTypes, fileType) {
			err2 := apierror.ErrInvalidFileType
			apierror.AddPayload(&err2, map[string]any{
				"allowed_types": allowedTypes,
				"received_type": fileType,
			})
			return err2
		}

        imageUrl, err = config.UploadFile("course/image/" +id.String()+ "." + imageFile.Filename, imageFile)
        if err != nil {
            return fmt.Errorf("failed to upload image: %v", err)
        }
    }
	

    // Upload syllabus if present
    if syllabusFile != nil {

        fileType, err := fileutil.DetectMultipartFileType(syllabusFile)
        
        if err != nil {
			log.Println("Error detecting syllabus type: ", err)
			return apierror.ErrInternalServer
		}
        allowedTypes := fileutil.SyllabusContentTypes
		if !slices.Contains(allowedTypes, fileType) {
			err2 := apierror.ErrInvalidFileType
			apierror.AddPayload(&err2, map[string]any{
				"allowed_types": allowedTypes,
				"received_type": fileType,
			})
			return err2
		}


        syllabusUrl, err = config.UploadFile("course/syllabus/" +id.String()+ "." +  syllabusFile.Filename, syllabusFile)
        if err != nil {
            return fmt.Errorf("failed to upload syllabus: %v", err)
        }
    }





    // Create course entity
    course := Course{
        Title:        req.Title,
        Description:  req.Description,
        Price:        req.Price,
        ImageURL:     imageUrl,
        SyllabusURL:  syllabusUrl,
        InstructorID: uuidInstructorID,
        Difficulty:   req.Difficulty,
        ID:           id,
    }

    return uc.courseRepo.Create(ctx, &course)
}

func (uc *UseCase) GetByInstructorID(ctx context.Context, instructorID uuid.UUID) ([]Course, error) {
    return uc.courseRepo.FindByInstructorID(ctx, instructorID)
}

func (uc *UseCase) Update(ctx context.Context, req UpdateCourseRequest, id uuid.UUID, imageFile, syllabusFile *multipart.FileHeader) (Course, error) {
    // Fetch the existing course to update
    course, err := uc.courseRepo.GetByID(ctx, id)
    if err != nil {
        return Course{}, fmt.Errorf("could not find course: %v", err)
    }

    // Update fields if provided
    if req.Title != nil {
        course.Title = *req.Title
    }
    if req.Description != nil {
        course.Description = *req.Description
    }
    if req.Price != nil {
        course.Price = *req.Price
    }
    if req.Difficulty != nil {
        course.Difficulty = *req.Difficulty
    }

    // Handle image update if file is provided
    if imageFile != nil {

        if imageFile.Size > 2*fileutil.MegaByte {
			err2 := apierror.ErrFileTooLarge
			apierror.AddPayload(&err2, map[string]string{
				"max_size":      "2 MB",
				"received_size": fileutil.ByteToAppropriateUnit(imageFile.Size),
			})
			return Course{},err2
		}
        fileType, err := fileutil.DetectMultipartFileType(imageFile)

        if err != nil {
			log.Println("Error detecting image type: ", err)
			return Course{},apierror.ErrInternalServer
		}

        allowedTypes := fileutil.ImageContentTypes
		if !slices.Contains(allowedTypes, fileType) {
			err2 := apierror.ErrInvalidFileType
			apierror.AddPayload(&err2, map[string]any{
				"allowed_types": allowedTypes,
				"received_type": fileType,
			})
			return Course{},err2
		}
        imageUrl, err := config.UploadFile("course/image/" +id.String()+ "." +imageFile.Filename, imageFile) // Assumes UploadFile encapsulates S3 logic
        if err != nil {
            return Course{}, fmt.Errorf("failed to upload image: %v", err)
        }
        course.ImageURL = imageUrl
    }

    // Handle syllabus update if file is provided
    if syllabusFile != nil {
        fileType, err := fileutil.DetectMultipartFileType(syllabusFile)
        
        if err != nil {
			log.Println("Error detecting syllabus type: ", err)
			return Course{},apierror.ErrInternalServer
		}
        allowedTypes := fileutil.SyllabusContentTypes
		if !slices.Contains(allowedTypes, fileType) {
			err2 := apierror.ErrInvalidFileType
			apierror.AddPayload(&err2, map[string]any{
				"allowed_types": allowedTypes,
				"received_type": fileType,
			})
			return Course{},err2
		}
        syllabusUrl, err := config.UploadFile("course/syllabus/" +id.String()+ "." +syllabusFile.Filename, syllabusFile)
        if err != nil {
            return Course{}, fmt.Errorf("failed to upload syllabus: %v", err)
        }
        course.SyllabusURL = syllabusUrl
    }

    // Update the course in the repository
    err = uc.courseRepo.Update(ctx, &course)
    if err != nil {
        return Course{}, fmt.Errorf("failed to update course: %v", err)
    }

    return course, nil
}


func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
    return uc.courseRepo.Delete(ctx, id)
}
