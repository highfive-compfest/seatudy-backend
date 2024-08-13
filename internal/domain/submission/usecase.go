package submission

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/assignment"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/attachment"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/course"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/courseenroll"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
)

type UseCase struct {
	repo              Repository
	assignmentRepo    assignment.Repository
	courseRepo        course.Repository
	attachmentUseCase attachment.UseCase
	courseEnrollRepo  courseenroll.Repository
}

// NewUseCase creates a new instance of the submission use case.
func NewUseCase(repo Repository, aRepo assignment.Repository, auc attachment.UseCase,courseRepo course.Repository, ceRepo courseenroll.Repository) *UseCase {
	return &UseCase{repo: repo, assignmentRepo: aRepo, attachmentUseCase: auc,courseRepo: courseRepo, courseEnrollRepo: ceRepo}
}

// CreateSubmission handles the business logic for creating a new submission.
func (uc *UseCase) CreateSubmission(ctx context.Context, req *CreateSubmissionRequest, userId string) error {

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return apierror.ErrInternalServer
	}

	assignmentUUID, err := uuid.Parse(req.AssignmentID)
	if err != nil {
		return apierror.ErrInternalServer
	}

	err = uc.CheckSubmissionExists(ctx,userUUID,assignmentUUID)
	if err != nil {
		return err
	}

	err = uc.VerifyCourseEnroll(ctx, userUUID, assignmentUUID)
	if err != nil {
		return err
	}

	id, err := uuid.NewV7()
	if err != nil {

		return apierror.ErrInternalServer
	}

	// check valid assignment ID
	assignment, err := uc.assignmentRepo.GetByID(ctx, assignmentUUID)
	if err != nil {
		return ErrAssignmentNotFound
	}

	submission := &schema.Submission{
		ID:           id,
		AssignmentID: assignment.ID,
		UserID:       userUUID,
		Content:      req.Content,
	}

	log.Println(req.Attachments)

	for _, fileHeader := range req.Attachments {
		attachment, err := uc.attachmentUseCase.CreateSubmissionAttachment(ctx, fileHeader, "")
		if err != nil {

			return ErrS3UploadFail
		}
		submission.Attachments = append(submission.Attachments, attachment)
	}

	if err := uc.repo.Create(ctx, submission); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GradeSubmission(ctx context.Context, userId string, id uuid.UUID, grade float64) error {
	submission, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	log.Println(submission)

	assignment, err := uc.assignmentRepo.GetByID(ctx, submission.AssignmentID)
	if err != nil {
		return err
	}

	log.Println(assignment.CourseID)

	courseId := assignment.CourseID
	// Get the course to verify the instructor ID
	course, err := uc.courseRepo.GetByID(ctx, courseId)
	if err != nil {
		return err
	}
	log.Println("halo guys")
	log.Println(course)

	// Check if the current user is the instructor of the course
	if course.InstructorID.String() != userId {
		return ErrNotOwnerCourse
	}

	submission.Grade = grade
	return uc.repo.Update(ctx, submission)
}

// UpdateSubmission handles the business logic for updating an existing submission.
func (uc *UseCase) UpdateSubmission(ctx context.Context, id uuid.UUID, req *UpdateSubmissionRequest, userId string) error {
	submission, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if submission.UserID.String() != userId {
		return ErrNotOwnerSubmission
	}

	if req.Content != nil {
		submission.Content = *req.Content
	}

	if req.Attachments != nil {
		for _, att := range submission.Attachments {
			err := uc.attachmentUseCase.DeleteAttachment(ctx, att.ID)
			if err != nil {
				return err
			}
		}

		submission.Attachments = []schema.Attachment{} 

		for _, fileHeader := range req.Attachments {
			attachment, err := uc.attachmentUseCase.CreateSubmissionAttachment(ctx, fileHeader, "")
			if err != nil {

				return ErrS3UploadFail
			}
			submission.Attachments = append(submission.Attachments, attachment)
		}
	}

	return uc.repo.Update(ctx, submission)
}

// DeleteSubmission handles the business logic for deleting a submission.
func (uc *UseCase) DeleteSubmission(ctx context.Context, id uuid.UUID, userId string) error {
	submission, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if submission.UserID.String() != userId {
		return ErrNotOwnerSubmission
	}

	return uc.repo.Delete(ctx, id)
}

// GetSubmissionByID handles fetching a submission by its ID.
func (uc *UseCase) GetSubmissionByID(ctx context.Context, id uuid.UUID) (*schema.Submission, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetAllSubmissionsByAssignment handles fetching all submissions for a given assignment.
func (uc *UseCase) GetAllSubmissionsByAssignment(ctx context.Context, assignmentID uuid.UUID) ([]schema.Submission, error) {
	return uc.repo.GetAllByAssignment(ctx, assignmentID)
}

func (uc *UseCase) VerifyCourseEnroll(ctx context.Context, userID uuid.UUID, assignmentID uuid.UUID) error {
	ass, err := uc.assignmentRepo.GetByID(ctx, assignmentID)
	if err != nil {
		return err
	}

	enroll, err := uc.courseEnrollRepo.IsEnrolled(ctx, userID, ass.CourseID)
	if err != nil {
		return apierror.ErrInternalServer
	}

	if !enroll {
		return ErrNotEnrollCourse
	}

	return nil
}

func (uc *UseCase) CheckSubmissionExists(ctx context.Context, userID uuid.UUID, assignmentID uuid.UUID) error {
	exists, err := uc.repo.CheckSubmissionExists(ctx, userID, assignmentID)
	if err != nil {
		return err
	}
	if exists {
		return ErrSubmissionAlreadyExists // Define this error in your errors package
	}
	return nil
}
