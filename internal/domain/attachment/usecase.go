// file: attachment/usecase.go

package attachment

import (
	"context"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"log"

	"mime/multipart"

	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/config"
)

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (auc *UseCase) CreateAttachment(ctx context.Context, fileHeader *multipart.FileHeader, description string,materialID uuid.UUID) (schema.Attachment, error) {
	// Validate file type
	// fileType, err := fileutil.DetectMultipartFileType(fileHeader)

	// if err != nil {
	// 	log.Println("Error detecting file type: ", err)
	// 	return uuid.Nil,apierror.ErrInternalServer
	// }

	// allowedTypes := fileutil.ImageContentTypes
	// if !slices.Contains(allowedTypes, fileType) {
	// 	err2 := apierror.ErrInvalidFileType
	// 	apierror.AddPayload(&err2, map[string]any{
	// 		"allowed_types": allowedTypes,
	// 		"received_type": fileType,
	// 	})
	// 	return uuid.Nil, err2
	// }

	id, err := uuid.NewV7()
	if err != nil {
		log.Println("Error generating UUID: ", err)
		return schema.Attachment{}, apierror.ErrInternalServer
	}

	log.Println("masuk ini")
	// Upload file and get URL
	fileURL, err := config.UploadFile("attachments/"+id.String()+"."+fileHeader.Filename, fileHeader)
	if err != nil {
		return schema.Attachment{}, err
	}

	// Create attachment record
	att := schema.Attachment{
		ID:          id,
		URL:         fileURL,
		Description: description,
		MaterialID: &materialID,
	}
	if err := auc.repo.Create(ctx, &att); err != nil {
		return schema.Attachment{}, apierror.ErrInternalServer
	}

	return att, nil
}

func (auc *UseCase) CreateAssignmentAttachment(ctx context.Context, fileHeader *multipart.FileHeader, description string,assignmentID uuid.UUID) (schema.Attachment, error) {

	id, err := uuid.NewV7()
	if err != nil {
		log.Println("Error generating UUID: ", err)
		return schema.Attachment{}, apierror.ErrInternalServer
	}

	log.Println("masuk ini")
	// Upload file and get URL
	fileURL, err := config.UploadFile("attachments/"+id.String()+"."+fileHeader.Filename, fileHeader)
	if err != nil {
		return schema.Attachment{}, err
	}

	// Create attachment record
	att := schema.Attachment{
		ID:          id,
		URL:         fileURL,
		Description: description,
		AssignmentID: &assignmentID,
	}
	if err := auc.repo.Create(ctx, &att); err != nil {
		return schema.Attachment{}, apierror.ErrInternalServer
	}

	return att, nil
}

func (uc *UseCase) GetAttachmentByID(ctx context.Context, id uuid.UUID) (*schema.Attachment, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) UpdateAttachment(ctx context.Context, id uuid.UUID, req AttachmentUpdateRequest) (*schema.Attachment, error) {
	attachment, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.File != nil {
		fileURL, err := config.UploadFile("attachments/"+id.String()+"."+req.File.Filename, req.File)
		if err != nil {
			return nil, ErrS3UploadFail
		}
		attachment.URL = fileURL
	}

	if req.Description != "" {
		attachment.Description = req.Description
	}

	if err := uc.repo.Update(ctx, attachment); err != nil {
		return nil, err
	}

	return attachment, nil
}

func (uc *UseCase) DeleteAttachment(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}
