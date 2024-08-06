package course

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"github.com/highfive-compfest/seatudy-backend/internal/config"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"

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


    // Upload image if present
    if imageFile != nil {
        imageUrl, err = config.UploadFile( imageFile.Filename, imageFile)
        if err != nil {
            return fmt.Errorf("failed to upload image: %v", err)
        }
    }
	

    // Upload syllabus if present
    if syllabusFile != nil {
        syllabusUrl, err = config.UploadFile( syllabusFile.Filename, syllabusFile)
        if err != nil {
            return fmt.Errorf("failed to upload syllabus: %v", err)
        }
    }

	id, err := uuid.NewV7()
	if err != nil {
		log.Println("Error generating UUID: ", err)
		return apierror.ErrInternalServer
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
        imageUrl, err := config.UploadFile(imageFile.Filename, imageFile) // Assumes UploadFile encapsulates S3 logic
        if err != nil {
            return Course{}, fmt.Errorf("failed to upload image: %v", err)
        }
        course.ImageURL = imageUrl
    }

    // Handle syllabus update if file is provided
    if syllabusFile != nil {
        syllabusUrl, err := config.UploadFile(syllabusFile.Filename, syllabusFile)
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
