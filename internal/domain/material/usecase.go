package material

import (
	"context"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"

	"log"

	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/attachment"
)

type UseCase struct {
	repo              Repository
	attachmentUseCase *attachment.UseCase // Add this line
}

func NewUseCase(repo Repository, attachmentUseCase *attachment.UseCase) *UseCase {
	return &UseCase{repo: repo, attachmentUseCase: attachmentUseCase}
}

func (uc *UseCase) CreateMaterial(ctx context.Context, req CreateMaterialRequest) error {
	// Create a new material instance

	courseId, err := uuid.Parse(req.CourseID)

	if err != nil {
		return apierror.ErrInternalServer

	}
	id, err := uuid.NewV7()
	if err != nil {
		log.Println("Error generating UUID: ", err)
		return apierror.ErrInternalServer
	}
	mat := schema.Material{
		ID:          id,
		CourseID:    courseId,
		Title:       req.Title,
		Description: req.Description,
	}

	log.Println("halo masuk attachemnt")

	log.Println(req.Attachments)

	for _, attachmentInput := range req.Attachments {
		attachment, err := uc.attachmentUseCase.CreateAttachment(ctx, attachmentInput.File, attachmentInput.Description)
		if err != nil {
			log.Print("failed")
			return ErrS3UploadFail
		}

		log.Println("keluuar attacment", attachment)
		// The attachment is now directly linked to the material by MaterialID
		mat.Attachments = append(mat.Attachments, attachment)
	}

	// Save the material with its attachments
	return uc.repo.Create(ctx, &mat)
}

func (uc *UseCase) GetMaterialByID(ctx context.Context, id uuid.UUID) (*schema.Material, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetAllMaterials(ctx context.Context) ([]*schema.Material, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *UseCase) UpdateMaterial(ctx context.Context, req UpdateMaterialRequest, id uuid.UUID) error {
	mat, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return ErrMaterialNotFound
	}

	// Update the material fields from the request
	if req.Title != nil {
		mat.Title = *req.Title
	}
	if req.Description != nil {
		mat.Description = *req.Description
	}

	// Handle attachments update
	for _, attachmentInput := range req.Attachments {
		if attachmentInput.File != nil {
			attachment, err := uc.attachmentUseCase.CreateAttachment(ctx, attachmentInput.File, attachmentInput.Description)
			if err != nil {

				return ErrS3UploadFail
			}
			mat.Attachments = append(mat.Attachments, attachment)
		}
	}
	return uc.repo.Update(ctx, mat)
}

func (uc *UseCase) DeleteMaterial(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}
