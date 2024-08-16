// file: attachment/usecase.go

package attachment

import (
	"context"
	"log"

	"github.com/highfive-compfest/seatudy-backend/internal/schema"

	"mime/multipart"

	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/config"
)

type UseCase struct {
	repo     Repository
	uploader config.FileUploader
}

func NewUseCase(repo Repository, uploader config.FileUploader) *UseCase {
	return &UseCase{repo: repo, uploader: uploader}
}

func (auc *UseCase) CreateAttachment(ctx context.Context, fileHeader *multipart.FileHeader, description string, materialID uuid.UUID) (schema.Attachment, error) {

	id, err := uuid.NewV7()
	if err != nil {
		log.Println("Error generating UUID: ", err)
		return schema.Attachment{}, apierror.ErrInternalServer.Build()
	}

	log.Println("masuk ini")

	fileURL, err := auc.uploader.UploadFile("attachments/material/"+id.String()+"."+fileHeader.Filename, fileHeader)
	if err != nil {
		return schema.Attachment{}, err
	}


	att := schema.Attachment{
		ID:          id,
		URL:         fileURL,
		Description: description,
		MaterialID:  &materialID,
	}
	if err := auc.repo.Create(ctx, &att); err != nil {
		return schema.Attachment{}, apierror.ErrInternalServer.Build()
	}

	return att, nil
}

func (auc *UseCase) CreateAssignmentAttachment(ctx context.Context, fileHeader *multipart.FileHeader, description string, assignmentID uuid.UUID) (schema.Attachment, error) {

	id, err := uuid.NewV7()
	if err != nil {
		log.Println("Error generating UUID: ", err)
		return schema.Attachment{}, apierror.ErrInternalServer.Build()
	}

	log.Println("masuk ini")

	fileURL, err := auc.uploader.UploadFile("attachments/assignment/"+id.String()+"."+fileHeader.Filename, fileHeader)
	if err != nil {
		return schema.Attachment{}, err
	}


	att := schema.Attachment{
		ID:           id,
		URL:          fileURL,
		Description:  description,
		AssignmentID: &assignmentID,
	}
	if err := auc.repo.Create(ctx, &att); err != nil {
		return schema.Attachment{}, apierror.ErrInternalServer.Build()
	}

	return att, nil
}

func (auc *UseCase) CreateSubmissionAttachment(ctx context.Context, fileHeader *multipart.FileHeader, description string) (schema.Attachment, error) {

	id, err := uuid.NewV7()
	if err != nil {
		log.Println("Error generating UUID: ", err)
		return schema.Attachment{}, apierror.ErrInternalServer.Build()
	}

	log.Println("masuk ini")

	fileURL, err := auc.uploader.UploadFile("attachments/submission/"+id.String()+"."+fileHeader.Filename, fileHeader)
	if err != nil {
		return schema.Attachment{}, err
	}


	att := schema.Attachment{
		ID:          id,
		URL:         fileURL,
		Description: description,
	}
	if err := auc.repo.Create(ctx, &att); err != nil {
		return schema.Attachment{}, apierror.ErrInternalServer.Build()
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
		fileURL, err := uc.uploader.UploadFile("attachments/"+id.String()+"."+req.File.Filename, req.File)
		if err != nil {
			return nil, ErrS3UploadFail.Build()
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
