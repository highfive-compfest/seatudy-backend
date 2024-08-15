package assignment

import (
	"context"

	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/attachment"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
)

type UseCase struct {
	repo              Repository
	attachmentUseCase *attachment.UseCase // Add this line
}

func NewUseCase(repo Repository, attachmentUseCase *attachment.UseCase) *UseCase {
	return &UseCase{repo: repo, attachmentUseCase: attachmentUseCase}
}

func (uc *UseCase) CreateAssignment(ctx context.Context, req CreateAssignmentRequest, courseId uuid.UUID) error {

	id, err := uuid.NewV7()
	if err != nil {
		return apierror.ErrInternalServer.Build()
	}
	assignment := &schema.Assignment{
		ID:          id,
		CourseID:    courseId,
		Title:       req.Title,
		Description: req.Description,
		Due:         req.Due,
	}
	return uc.repo.Create(ctx, assignment)
}

func (uc *UseCase) UpdateAssignment(ctx context.Context, id uuid.UUID, req UpdateAssignmentRequest) error {
	assignment, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Title != nil {
		assignment.Title = *req.Title
	}
	if req.Description != nil {
		assignment.Description = *req.Description
	}
	if req.Due != nil {
		assignment.Due = req.Due
	}

	return uc.repo.Update(ctx, assignment)
}

func (uc *UseCase) DeleteAssignment(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *UseCase) GetAssignmentByID(ctx context.Context, id uuid.UUID) (*schema.Assignment, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetAssignmentsByCourse(ctx context.Context, courseId uuid.UUID) ([]*schema.Assignment, error) {
	return uc.repo.GetByCourseID(ctx, courseId)
}

func (uc *UseCase) AddAttachment(ctx context.Context, id uuid.UUID, req AttachmentInput) error {
	ass, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return ErrAssignmentNotFound.Build()
	}

	if req.File != nil {

		attachment, err := uc.attachmentUseCase.CreateAssignmentAttachment(ctx, req.File, req.Description, id)
		if err != nil {
			return ErrS3UploadFail.Build()
		}
		ass.Attachments = append(ass.Attachments, attachment)
	}

	return uc.repo.Update(ctx, ass)
}
