package assignment

import (
	"context"
	"mime/multipart"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/attachment"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, a *schema.Assignment) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, a *schema.Assignment) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*schema.Assignment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*schema.Assignment), args.Error(1)
}

func (m *MockRepository) GetByCourseID(ctx context.Context, courseId uuid.UUID) ([]*schema.Assignment, error) {
	args := m.Called(ctx, courseId)
	return args.Get(0).([]*schema.Assignment), args.Error(1)
}

type MockAttachmentRepository struct {
	mock.Mock
}

func (m *MockAttachmentRepository) Create(ctx context.Context, att *schema.Attachment) error {
	args := m.Called(ctx, att)
	return args.Error(0)
}

func (m *MockAttachmentRepository) Update(ctx context.Context, att *schema.Attachment) error {
	args := m.Called(ctx, att)
	return args.Error(0)
}

func (m *MockAttachmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*schema.Attachment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*schema.Attachment), args.Error(1)
}

func (m *MockAttachmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockFileUploader struct {
	mock.Mock
}

func (m *MockFileUploader) UploadFile(key string, fileHeader *multipart.FileHeader) (string, error) {
	args := m.Called(key, fileHeader)
	return args.String(0), args.Error(1)
}

type AssignmentUseCaseTestSuite struct {
	suite.Suite
	attachmentRepo    *MockAttachmentRepository
	uploader          *MockFileUploader
	assignmentRepo    *MockRepository
	assignmentUseCase *UseCase
	attachmentUseCase *attachment.UseCase
}

func (suite *AssignmentUseCaseTestSuite) SetupTest() {
	suite.attachmentRepo = new(MockAttachmentRepository)
	suite.uploader = new(MockFileUploader)
	suite.assignmentRepo = new(MockRepository)
	suite.attachmentUseCase = attachment.NewUseCase(suite.attachmentRepo, suite.uploader)
	suite.assignmentUseCase = NewUseCase(suite.assignmentRepo, suite.attachmentUseCase)

}

func (suite *AssignmentUseCaseTestSuite) TestCreateAssignment_Success() {
	ctx := context.Background()
	courseId, _ := uuid.NewV7()
	newDueTime := time.Now().Add(72 * time.Hour)
	req := CreateAssignmentRequest{
		Title:       "New Assignment",
		Description: "Description of the new assignment",

		Due: &newDueTime,
	}

	suite.assignmentRepo.On("Create", ctx, mock.AnythingOfType("*schema.Assignment")).Return(nil)

	err := suite.assignmentUseCase.CreateAssignment(ctx, req, courseId)

	assert.NoError(suite.T(), err)
	suite.assignmentRepo.AssertExpectations(suite.T())
}

func (suite *AssignmentUseCaseTestSuite) TestUpdateAssignment_Success() {
	ctx := context.Background()
	id, _ := uuid.NewV7()
	newTitle := "Updated Title"
	newDescription := "Updated Description"
	newDueTime := time.Now().Add(72 * time.Hour) // Create a time.Time variable

	req := UpdateAssignmentRequest{
		Title:       &newTitle,
		Description: &newDescription,
		Due:         &newDueTime, // Pass the address of newDueTime
	}

	beforeDueTime := time.Now().Add(24 * time.Hour)

	existingAssignment := &schema.Assignment{
		ID:          id,
		Title:       "Old Title",
		Description: "Old Description",
		Due:         &beforeDueTime,
	}

	suite.assignmentRepo.On("GetByID", ctx, id).Return(existingAssignment, nil)
	suite.assignmentRepo.On("Update", ctx, mock.AnythingOfType("*schema.Assignment")).Return(nil)

	err := suite.assignmentUseCase.UpdateAssignment(ctx, id, req)

	assert.NoError(suite.T(), err)
	suite.assignmentRepo.AssertExpectations(suite.T())
}

func (suite *AssignmentUseCaseTestSuite) TestDeleteAssignment_Success() {
	ctx := context.Background()
	id := uuid.New()
	suite.assignmentRepo.On("Delete", ctx, id).Return(nil)
	err := suite.assignmentUseCase.DeleteAssignment(ctx, id)
	assert.NoError(suite.T(), err)
	suite.assignmentRepo.AssertExpectations(suite.T())
}

func (suite *AssignmentUseCaseTestSuite) TestGetAssignmentByID_Success() {
	ctx := context.Background()
	id := uuid.New()
	expectedAssignment := &schema.Assignment{ID: id}
	suite.assignmentRepo.On("GetByID", ctx, id).Return(expectedAssignment, nil)
	assignment, err := suite.assignmentUseCase.GetAssignmentByID(ctx, id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedAssignment, assignment)
	suite.assignmentRepo.AssertExpectations(suite.T())
}

func (suite *AssignmentUseCaseTestSuite) TestGetAssignmentsByCourse_Success() {
	ctx := context.Background()
	courseId := uuid.New()
	expectedAssignments := []*schema.Assignment{{ID: uuid.New()}, {ID: uuid.New()}}
	suite.assignmentRepo.On("GetByCourseID", ctx, courseId).Return(expectedAssignments, nil)
	assignments, err := suite.assignmentUseCase.GetAssignmentsByCourse(ctx, courseId)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedAssignments, assignments)
	suite.assignmentRepo.AssertExpectations(suite.T())
}

func (suite *AssignmentUseCaseTestSuite) TestAddAttachment_Success() {
	ctx := context.Background()
	id := uuid.New()
	fileHeader := &multipart.FileHeader{Filename: "assignment.pdf", Size: 1 * 1024 * 1024}
	req := AttachmentInput{File: fileHeader, Description: "Assignment Attachment"}

	existingAssignment := &schema.Assignment{ID: id, Attachments: []schema.Attachment{}}
	suite.assignmentRepo.On("GetByID", ctx, id).Return(existingAssignment, nil)
	suite.uploader.On("UploadFile", mock.Anything, mock.Anything).Return("http://example.com/assignment.pdf", nil)
	suite.attachmentRepo.On("Create", ctx, mock.Anything).Return(nil)
	suite.assignmentRepo.On("Update", ctx, existingAssignment).Return(nil)

	err := suite.assignmentUseCase.AddAttachment(ctx, id, req)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), existingAssignment.Attachments, 1)
	assert.Equal(suite.T(), "http://example.com/assignment.pdf", existingAssignment.Attachments[0].URL)
	assert.Equal(suite.T(), "Assignment Attachment", existingAssignment.Attachments[0].Description)
	suite.assignmentRepo.AssertExpectations(suite.T())
	suite.attachmentRepo.AssertExpectations(suite.T())
}

func TestCourseUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(AssignmentUseCaseTestSuite))
}
