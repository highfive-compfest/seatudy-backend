package submission

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/config"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/assignment"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/attachment"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/course"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/courseenroll"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/notification"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/user"
	"github.com/highfive-compfest/seatudy-backend/internal/mailer"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"log"
)

type UseCase struct {
	repo              Repository
	assignmentRepo    assignment.Repository
	courseRepo        course.Repository
	attachmentUseCase attachment.UseCase
	courseEnrollRepo  courseenroll.Repository
	userRepo          user.IRepository
	notifRepo         notification.IRepository
	mailDialer        config.IMailer
}

// NewUseCase creates a new instance of the submission use case.
func NewUseCase(repo Repository, aRepo assignment.Repository, auc attachment.UseCase, courseRepo course.Repository,
	ceRepo courseenroll.Repository, userRepo user.IRepository, notifRepo notification.IRepository, mailDialer config.IMailer) *UseCase {
	return &UseCase{repo: repo, assignmentRepo: aRepo, attachmentUseCase: auc, courseRepo: courseRepo,
		courseEnrollRepo: ceRepo, userRepo: userRepo, notifRepo: notifRepo, mailDialer: mailDialer}
}

//go:embed new_submission_instructor_email_template.html
var newSubmissionInstructorEmailTemplate string

// CreateSubmission handles the business logic for creating a new submission.
func (uc *UseCase) CreateSubmission(ctx context.Context, req *CreateSubmissionRequest, userId string) error {

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return apierror.ErrInternalServer.Build()
	}

	assignmentUUID, err := uuid.Parse(req.AssignmentID)
	if err != nil {
		return apierror.ErrInternalServer.Build()
	}

	err = uc.CheckSubmissionExists(ctx, userUUID, assignmentUUID)
	if err != nil {
		return err
	}

	err = uc.VerifyCourseEnroll(ctx, userUUID, assignmentUUID)
	if err != nil {
		return err
	}

	id, err := uuid.NewV7()
	if err != nil {

		return apierror.ErrInternalServer.Build()
	}

	// check valid assignment ID
	assignmentObj, err := uc.assignmentRepo.GetByID(ctx, assignmentUUID)
	if err != nil {
		return ErrAssignmentNotFound
	}

	submission := &schema.Submission{
		ID:           id,
		AssignmentID: assignmentObj.ID,
		UserID:       userUUID,
		Content:      req.Content,
	}

	log.Println(req.Attachments)

	for _, fileHeader := range req.Attachments {
		attachmentObj, err := uc.attachmentUseCase.CreateSubmissionAttachment(ctx, fileHeader, "")
		if err != nil {

			return ErrS3UploadFail
		}
		submission.Attachments = append(submission.Attachments, attachmentObj)
	}

	if err := uc.repo.Create(ctx, submission); err != nil {
		return err
	}

	go func() {
		studentName := ctx.Value("user.name").(string)
		studentEmail := ctx.Value("user.email").(string)

		courseObj, err := uc.courseRepo.GetByID(ctx, assignmentObj.CourseID)
		if err != nil {
			log.Println("Error getting course: ", err)
			return
		}

		instructor, err := uc.userRepo.GetByID(courseObj.InstructorID)
		if err != nil {
			log.Println("Error getting instructor: ", err)
			return
		}

		// Send email to instructor
		go func() {
			emailData := map[string]any{
				"instructor_name":  instructor.Name,
				"course_title":     courseObj.Title,
				"assignment_title": assignmentObj.Title,
				"student_name":     studentName,
				"student_email":    studentEmail,
			}

			mail, err := mailer.GenerateMail(instructor.Email, "New Assignment Submission", newSubmissionInstructorEmailTemplate, emailData)
			if err != nil {
				log.Println("Error generating email: ", err)
				return
			}

			if err = uc.mailDialer.DialAndSend(mail); err != nil {
				log.Println("Error sending email: ", err)
				return
			}
		}()

		// Create in-app notification
		go func() {
			notificationID, err := uuid.NewV7()
			if err != nil {
				log.Println("Error generating notification ID: ", err)
				return
			}
			notif := schema.Notification{
				ID:     notificationID,
				UserID: courseObj.InstructorID,
				Title:  "New Assignment Submission",
				Detail: fmt.Sprintf("%s has submitted an assignment for %s", studentName, assignmentObj.Title),
			}

			if err := uc.notifRepo.Create(&notif); err != nil {
				log.Println("Error creating notification: ", err)
				return
			}
		}()
	}()

	return nil
}

//go:embed submission_graded_student_email_template.html
var submissionGradedStudentEmailTemplate string

func (uc *UseCase) GradeSubmission(ctx context.Context, userId string, id uuid.UUID, grade float64) error {
	submission, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	assignmentObj, err := uc.assignmentRepo.GetByID(ctx, submission.AssignmentID)
	if err != nil {
		return err
	}

	courseId := assignmentObj.CourseID
	// Get the course to verify the instructor ID
	courseObj, err := uc.courseRepo.GetByID(ctx, courseId)
	if err != nil {
		return err
	}

	// Check if the current user is the instructor of the course
	if courseObj.InstructorID.String() != userId {
		return ErrNotOwnerCourse
	}

	submission.Grade = grade
	if err := uc.repo.Update(ctx, submission); err != nil {
		log.Println("Error updating submission: ", err)
		return err
	}

	// Send email to student
	go func() {
		student, err := uc.userRepo.GetByID(submission.UserID)
		if err != nil {
			log.Println("Error getting student: ", err)
			return
		}

		emailData := map[string]any{
			"student_name":     student.Name,
			"course_title":     courseObj.Title,
			"assignment_title": assignmentObj.Title,
			"grade":            grade,
		}

		mail, err := mailer.GenerateMail(student.Email, "Assignment Graded",
			submissionGradedStudentEmailTemplate, emailData)
		if err != nil {
			log.Println("Error generating email: ", err)
			return
		}

		if err = uc.mailDialer.DialAndSend(mail); err != nil {
			log.Println("Error sending email: ", err)
			return
		}
	}()

	// Send in-app notification to student
	go func() {
		notifID, err := uuid.NewV7()
		if err != nil {
			return
		}

		notif := schema.Notification{
			ID:     notifID,
			UserID: submission.UserID,
			Title:  "Assignment Graded",
			Detail: fmt.Sprintf("Your submission for %s in course %s has been graded", assignmentObj.Title, courseObj.Title),
		}

		if err := uc.notifRepo.Create(&notif); err != nil {
			log.Println("Error creating notification: ", err)
			return
		}
	}()

	return nil
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
		return apierror.ErrInternalServer.Build()
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
