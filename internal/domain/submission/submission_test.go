package submission

import (
	"context"
	"mime/multipart"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/config"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/attachment"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/courseenroll"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gomail.v2"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, s *schema.Submission) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, s *schema.Submission) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*schema.Submission, error) {
	args := m.Called(ctx, id)
	if item := args.Get(0); item != nil {
		return item.(*schema.Submission), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) GetAllByAssignment(ctx context.Context, assignmentID uuid.UUID) ([]schema.Submission, error) {
	args := m.Called(ctx, assignmentID)
	if item := args.Get(0); item != nil {
		return item.([]schema.Submission), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) CheckSubmissionExists(ctx context.Context, userID, assignmentID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID, assignmentID)
	return args.Bool(0), args.Error(1)
}

type MockCourseRepository struct {
	mock.Mock
}

func (m *MockCourseRepository) GetAll(ctx context.Context, page, pageSize int) ([]schema.Course, int, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]schema.Course), args.Int(1), args.Error(2)
}

func (m *MockCourseRepository) GetByID(ctx context.Context, id uuid.UUID) (schema.Course, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(schema.Course), args.Error(1)
}

func (m *MockCourseRepository) GetRating(ctx context.Context, courseID uuid.UUID) (float32, int64, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).(float32), args.Get(1).(int64), args.Error(2)
}

func (m *MockCourseRepository) Create(ctx context.Context, course *schema.Course) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockCourseRepository) Update(ctx context.Context, course *schema.Course) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockCourseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCourseRepository) FindByInstructorID(ctx context.Context, instructorID uuid.UUID, page, pageSize int) ([]schema.Course, int, error) {
	args := m.Called(ctx, instructorID, page, pageSize)
	return args.Get(0).([]schema.Course), args.Int(1), args.Error(2)
}

func (m *MockCourseRepository) FindByPopularity(ctx context.Context, page, pageSize int) ([]schema.Course, int, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]schema.Course), args.Int(1), args.Error(2)
}

func (m *MockCourseRepository) GetUserCourseProgress(ctx context.Context, courseID, userID uuid.UUID) (float64, error) {
	args := m.Called(ctx, courseID, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockCourseRepository) SearchByTitle(ctx context.Context, title string, page, pageSize int) ([]schema.Course, int, error) {
	args := m.Called(ctx, title, page, pageSize)
	return args.Get(0).([]schema.Course), args.Int(1), args.Error(2)
}

func (m *MockCourseRepository) DynamicFilterCourses(ctx context.Context, filterType, filterValue, sort string, page, limit int) ([]schema.Course, int, error) {
	args := m.Called(ctx, filterType, filterValue, sort, page, limit)
	return args.Get(0).([]schema.Course), args.Int(1), args.Error(2)
}

type MockFileUploader struct {
	mock.Mock
}

func (m *MockFileUploader) UploadFile(key string, fileHeader *multipart.FileHeader) (string, error) {
	args := m.Called(key, fileHeader)
	return args.String(0), args.Error(1)
}

type MockEnrollRepository struct {
	mock.Mock
}

func (m *MockEnrollRepository) Create(ctx context.Context, enroll *schema.CourseEnroll) error {
	args := m.Called(ctx, enroll)
	return args.Error(0)
}

func (m *MockEnrollRepository) GetUsersByCourseID(ctx context.Context, courseID uuid.UUID) ([]schema.User, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).([]schema.User), args.Error(1)
}

func (m *MockEnrollRepository) GetCoursesByUserID(ctx context.Context, userID uuid.UUID) ([]schema.Course, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]schema.Course), args.Error(1)
}

func (m *MockEnrollRepository) IsEnrolled(ctx context.Context, userID, courseID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID, courseID)
	return args.Bool(0), args.Error(1)
}

type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) Create(notification *schema.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func (m *MockNotificationRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]*schema.Notification, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]*schema.Notification), args.Get(1).(int64), args.Error(2)
}

func (m *MockNotificationRepository) GetUnreadCount(userID uuid.UUID) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNotificationRepository) UpdateRead(notificationID uuid.UUID) error {
	args := m.Called(notificationID)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *schema.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*schema.User, error) {
	args := m.Called(id)
	user, ok := args.Get(0).(*schema.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*schema.User, error) {
	args := m.Called(email)
	user, ok := args.Get(0).(*schema.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) Update(user *schema.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateByEmail(email string, user *schema.User) error {
	args := m.Called(email, user)
	return args.Error(0)
}

type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) DialAndSend(msgs ...*gomail.Message) error {
	args := m.Called(msgs)
	return args.Error(0)
}

type MockAssignmentRepo struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockAssignmentRepo) Create(ctx context.Context, a *schema.Assignment) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

// Update mocks the Update method
func (m *MockAssignmentRepo) Update(ctx context.Context, a *schema.Assignment) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockAssignmentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetByID mocks the GetByID method
func (m *MockAssignmentRepo) GetByID(ctx context.Context, id uuid.UUID) (*schema.Assignment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*schema.Assignment), args.Error(1)
}

// GetByCourseID mocks the GetByCourseID method
func (m *MockAssignmentRepo) GetByCourseID(ctx context.Context, courseId uuid.UUID) ([]*schema.Assignment, error) {
	args := m.Called(ctx, courseId)
	return args.Get(0).([]*schema.Assignment), args.Error(1)
}

type MockAttachmentRepo struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockAttachmentRepo) Create(ctx context.Context, att *schema.Attachment) error {
	args := m.Called(ctx, att)
	return args.Error(0)
}

// Update mocks the Update method
func (m *MockAttachmentRepo) Update(ctx context.Context, att *schema.Attachment) error {
	args := m.Called(ctx, att)
	return args.Error(0)
}

// GetByID mocks the GetByID method
func (m *MockAttachmentRepo) GetByID(ctx context.Context, id uuid.UUID) (*schema.Attachment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*schema.Attachment), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockAttachmentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type SubmissionUseCaseTestSuite struct {
	suite.Suite

	courseRepo    *MockCourseRepository
	enrollRepo    *MockEnrollRepository
	assignmentRepo *MockAssignmentRepo
	attachmentRepo *MockAttachmentRepo
	userRepo *MockUserRepository
	mailer *MockMailer
	enrollUseCase *courseenroll.UseCase
	attachmentUseCase *attachment.UseCase
	notificationRepo *MockNotificationRepository
	submissionRepo *MockRepository
	submisionUseCase *UseCase
    uploader *MockFileUploader	
}

func (suite *SubmissionUseCaseTestSuite) SetupTest() {
	suite.submissionRepo = new(MockRepository)
	suite.assignmentRepo = new(MockAssignmentRepo)
	suite.courseRepo = new(MockCourseRepository)
	suite.enrollRepo = new(MockEnrollRepository)
	suite.attachmentRepo = new(MockAttachmentRepo)
	suite.userRepo =  new(MockUserRepository)
    suite.uploader = new(MockFileUploader)
	suite.mailer = new(MockMailer)
	suite.notificationRepo = new(MockNotificationRepository)
	suite.enrollUseCase = courseenroll.NewUseCase(suite.enrollRepo)
	suite.attachmentUseCase = attachment.NewUseCase(suite.attachmentRepo,suite.uploader)
	suite.submisionUseCase = NewUseCase(suite.submissionRepo,suite.assignmentRepo,*suite.attachmentUseCase,suite.courseRepo,suite.enrollRepo,suite.userRepo,suite.notificationRepo,suite.mailer)

}

func (suite *SubmissionUseCaseTestSuite) TestGradeSubmission_Success() {
	ctx := context.Background()
	userId,_ := uuid.NewV7()// This user is not the instructor
	submissionId,_ := uuid.NewV7()
	courseId,_ := uuid.NewV7()
	assignmentId,_ := uuid.NewV7()
	grade := 85.0

	submission := &schema.Submission{
		ID:           submissionId,
		AssignmentID: assignmentId,
		UserID:       uuid.New(),
	}

	assignment := &schema.Assignment{
		ID:       assignmentId,
		CourseID: courseId,
	}

	course := &schema.Course{
		ID:           courseId,
		InstructorID: userId,
	}

	suite.submissionRepo.On("GetByID", ctx, mock.Anything).Return(submission, nil)
	suite.assignmentRepo.On("GetByID", ctx, mock.Anything).Return(assignment, nil)
	suite.courseRepo.On("GetByID", ctx, mock.Anything).Return(*course, nil)
	suite.submissionRepo.On("Update", ctx, submission).Return(nil)

	// Call the function under test
	err := suite.submisionUseCase.GradeSubmission(ctx, userId.String(), submissionId, grade)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), grade, submission.Grade)
	suite.submissionRepo.AssertExpectations(suite.T())
	suite.assignmentRepo.AssertExpectations(suite.T())
	suite.courseRepo.AssertExpectations(suite.T())
}


func (suite *SubmissionUseCaseTestSuite) TestGradeSubmission_Failure_NotCourseInstructor() {
	ctx := context.Background()
	userId,_ := uuid.NewV7()// This user is not the instructor
	submissionId,_ := uuid.NewV7()
	courseId,_ := uuid.NewV7()
	assignmentId,_ := uuid.NewV7()

	submission := &schema.Submission{
		ID:           submissionId,
		AssignmentID: assignmentId,
		UserID:       uuid.New(),
	}

	assignment := &schema.Assignment{
		ID:       assignmentId,
		CourseID: courseId,
	}

	course := &schema.Course{
		ID:           courseId,
		InstructorID: uuid.New(), // Different instructor ID
	}

	suite.submissionRepo.On("GetByID", ctx, submissionId).Return(submission, nil)
	suite.assignmentRepo.On("GetByID", ctx, assignmentId).Return(assignment, nil)
	suite.courseRepo.On("GetByID", ctx, courseId).Return(*course, nil)

	// Call the function under test
	err := suite.submisionUseCase.GradeSubmission(ctx, userId.String(), submissionId, 90.0)

	// Assertions
	assert.Error(suite.T(), err)
	assert.NotEqual(suite.T(), 90.0, submission.Grade)
	assert.Equal(suite.T(), ErrNotOwnerCourse, err)
	suite.submissionRepo.AssertExpectations(suite.T())
	suite.assignmentRepo.AssertExpectations(suite.T())
	suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *SubmissionUseCaseTestSuite) TestDeleteSubmission_Success() {
	ctx := context.Background()
	id := uuid.New()
	userId := uuid.New().String()

	submission := &schema.Submission{
		ID:     id,
		UserID: uuid.MustParse(userId),
	}

	suite.submissionRepo.On("GetByID", ctx, id).Return(submission, nil)
	suite.submissionRepo.On("Delete", ctx, id).Return(nil)

	err := suite.submisionUseCase.DeleteSubmission(ctx, id, userId)
	assert.NoError(suite.T(), err)
	suite.submissionRepo.AssertExpectations(suite.T())
}

func (suite *SubmissionUseCaseTestSuite) TestDeleteSubmission_NotOwner() {
	ctx := context.Background()
	id := uuid.New()
	userId := uuid.New().String()
	otherUserId := uuid.New().String()

	submission := &schema.Submission{
		ID:     id,
		UserID: uuid.MustParse(otherUserId),
	}

	suite.submissionRepo.On("GetByID", ctx, id).Return(submission, nil)

	err := suite.submisionUseCase.DeleteSubmission(ctx, id, userId)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrNotOwnerSubmission, err)
	suite.submissionRepo.AssertExpectations(suite.T())
}

func (suite *SubmissionUseCaseTestSuite) TestGetSubmissionByID_Success() {
	ctx := context.Background()
	id := uuid.New()

	submission := &schema.Submission{
		ID: id,
	}

	suite.submissionRepo.On("GetByID", ctx, id).Return(submission, nil)

	result, err := suite.submisionUseCase.GetSubmissionByID(ctx, id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), submission, result)
	suite.submissionRepo.AssertExpectations(suite.T())
}

func (suite *SubmissionUseCaseTestSuite) TestGetAllSubmissionsByAssignment_Success() {
	ctx := context.Background()
	assignmentID := uuid.New()

	submissions := []schema.Submission{
		{ID: uuid.New()},
		{ID: uuid.New()},
	}

	suite.submissionRepo.On("GetAllByAssignment", ctx, assignmentID).Return(submissions, nil)

	result, err := suite.submisionUseCase.GetAllSubmissionsByAssignment(ctx, assignmentID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), submissions, result)
	suite.submissionRepo.AssertExpectations(suite.T())
}

func (suite *SubmissionUseCaseTestSuite) TestVerifyCourseEnroll_NotEnrolled() {
	ctx := context.Background()
	userID := uuid.New()
	assignmentID := uuid.New()
	courseID := uuid.New()

	assignment := &schema.Assignment{
		ID:       assignmentID,
		CourseID: courseID,
	}

	suite.assignmentRepo.On("GetByID", ctx, assignmentID).Return(assignment, nil)
	suite.enrollRepo.On("IsEnrolled", ctx, userID, courseID).Return(false, nil)

	err := suite.submisionUseCase.VerifyCourseEnroll(ctx, userID, assignmentID)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrNotEnrollCourse, err)
	suite.assignmentRepo.AssertExpectations(suite.T())
	suite.enrollRepo.AssertExpectations(suite.T())
}

func (suite *SubmissionUseCaseTestSuite) TestCheckSubmissionExists_Exists() {
	ctx := context.Background()
	userID := uuid.New()
	assignmentID := uuid.New()

	suite.submissionRepo.On("CheckSubmissionExists", ctx, userID, assignmentID).Return(true, nil)

	err := suite.submisionUseCase.CheckSubmissionExists(ctx, userID, assignmentID)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrSubmissionAlreadyExists, err)
	suite.submissionRepo.AssertExpectations(suite.T())
}

func (suite *SubmissionUseCaseTestSuite) TestCreateSubmission_Success() {
	ctx := context.WithValue(context.Background(), "user.name", "John Doe")
	ctx = context.WithValue(ctx, "user.email", "john.doe@example.com")
	_ = os.Setenv("ENV", "test")
	_ = os.Setenv("SMTP_EMAIL", "test")
	config.LoadEnv()
	userID := uuid.New().String()
	assignmentID := uuid.New().String()
	content := "Submission content here"
	fileHeader := multipart.FileHeader{Filename: "testfile.pdf", Size: 1024} // Mock file header
	req := &CreateSubmissionRequest{
		AssignmentID: assignmentID,
		Content:      content,
		Attachments:  []*multipart.FileHeader{&fileHeader},
	}

	userUUID, _ := uuid.Parse(userID)
	assignmentUUID, _ := uuid.Parse(assignmentID)
	assignment := &schema.Assignment{
		ID:       assignmentUUID,
		Title: "assignment",
		CourseID: uuid.New(),
	}
	instructorID := uuid.New()
	instructor := &schema.User{
		ID:   instructorID,
		Name: "Jane Instructor",
		Email: "jane.instructor@example.com",
	}

	suite.submissionRepo.On("CheckSubmissionExists", ctx, userUUID, assignmentUUID).Return(false, nil)
	suite.enrollRepo.On("IsEnrolled", ctx, userUUID, assignment.CourseID).Return(true, nil)
	suite.assignmentRepo.On("GetByID", ctx, assignmentUUID).Return(assignment, nil)
	suite.uploader.On("UploadFile", mock.Anything, mock.Anything).Return("http://example.com/testfile.pdf", nil)
	suite.attachmentRepo.On("Create", ctx, mock.Anything).Return(nil)
	suite.submissionRepo.On("Create", ctx, mock.Anything).Return(nil)
	suite.courseRepo.On("GetByID", ctx,mock.Anything).Return(schema.Course{ID: assignment.CourseID,Title: "mantap",Price: 10000, InstructorID: instructorID}, nil)
	suite.userRepo.On("GetByID", instructorID).Return(instructor, nil)
	suite.mailer.On("DialAndSend", mock.Anything).Return(nil)
	suite.notificationRepo.On("Create", mock.AnythingOfType("*schema.Notification")).Return(nil)

	assert.NotPanics(suite.T(), func() {
		err := suite.submisionUseCase.CreateSubmission(ctx, req, userID)
		assert.NoError(suite.T(), err)
	})

	suite.submissionRepo.AssertExpectations(suite.T())
	suite.assignmentRepo.AssertExpectations(suite.T())
	suite.enrollRepo.AssertExpectations(suite.T())
	suite.attachmentRepo.AssertExpectations(suite.T())

}

func (suite *SubmissionUseCaseTestSuite) TestUpdateSubmission_Success() {
    ctx := context.Background()
    userID := uuid.New()
    submissionID := uuid.New()
    submission := &schema.Submission{
        ID:           submissionID,
        UserID:       userID,
        Content:      "Original Content",
        Attachments:  []schema.Attachment{{ID: uuid.New(), Description: "Original"}},
    }

    newContent := "Updated Content"
    newFileHeader := multipart.FileHeader{Filename: "new_attachment.pdf", Size: 1024}
    newAttachmentID := uuid.New()
    newAttachment := schema.Attachment{ID: newAttachmentID, Description: "Updated", URL: "http://example.com/testfile.pdf"}

    suite.submissionRepo.On("GetByID", ctx, submissionID).Return(submission, nil)
    suite.submissionRepo.On("Update", ctx, submission).Return(nil)
    suite.attachmentRepo.On("Delete", ctx, mock.AnythingOfType("uuid.UUID")).Return(nil)
    suite.attachmentRepo.On("Create", ctx, mock.Anything).Run(func(args mock.Arguments) {
        arg := args.Get(1).(*schema.Attachment)
        *arg = newAttachment // directly modify the argument to simulate the creation
    }).Return(nil)
    suite.uploader.On("UploadFile", mock.Anything, mock.Anything).Return("http://example.com/testfile.pdf", nil)

    req := &UpdateSubmissionRequest{
        Content:     &newContent,
        Attachments: []*multipart.FileHeader{&newFileHeader},
    }

    err := suite.submisionUseCase.UpdateSubmission(ctx, submissionID, req, userID.String())

    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), newContent, submission.Content)
    assert.Len(suite.T(), submission.Attachments, 1)
    assert.Equal(suite.T(), newAttachment, submission.Attachments[0]) // This comparison will now work as expected

    suite.attachmentRepo.AssertExpectations(suite.T())
    suite.submissionRepo.AssertExpectations(suite.T())
}
func (suite *SubmissionUseCaseTestSuite) TestUpdateSubmission_Failure_NotOwner() {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	submissionID := uuid.New()
	submission := &schema.Submission{
		ID:           submissionID,
		UserID:       userID,
		Content:      "Original Content",
	}

	suite.submissionRepo.On("GetByID", ctx, submissionID).Return(submission, nil)

	req := &UpdateSubmissionRequest{
		Content:     nil, // No update to content
		Attachments: nil, // No new attachments
	}

	err := suite.submisionUseCase.UpdateSubmission(ctx, submissionID, req, otherUserID.String())

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrNotOwnerSubmission, err)

	suite.submissionRepo.AssertExpectations(suite.T())
}





func TestSubmissionUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(SubmissionUseCaseTestSuite))
}

