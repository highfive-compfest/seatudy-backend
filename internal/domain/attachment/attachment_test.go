package attachment

import (
	"context"
	"mime/multipart"
	"testing"

	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockRepository struct {
	mock.Mock
}


func (m *MockRepository) Create(ctx context.Context, att *schema.Attachment) error {
	args := m.Called(ctx, att)
	return args.Error(0)
}


func (m *MockRepository) Update(ctx context.Context, att *schema.Attachment) error {
	args := m.Called(ctx, att)
	return args.Error(0)
}


func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*schema.Attachment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*schema.Attachment), args.Error(1)
}


func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
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

type AttachmentUseCaseTestSuite struct {
	suite.Suite
	attachmentRepo *MockRepository
	attachmentUseCase * UseCase
	uploader *MockFileUploader
}

func (suite *AttachmentUseCaseTestSuite) SetupTest() {
	suite.attachmentRepo = new(MockRepository)
	suite.uploader =  new(MockFileUploader)
	suite.attachmentUseCase = NewUseCase(suite.attachmentRepo,suite.uploader )

}

func (suite *AttachmentUseCaseTestSuite) TestCreateAttachment_Success() {
	ctx := context.Background()
	fileHeader := &multipart.FileHeader{Filename: "test.pdf", Size: 1 * 1024 * 1024} // Assuming 1MB size
	description := "Test Description"
	materialID,_ := uuid.NewV7()

	expectedURL := "http://example.com/test.pdf"
	attachment := schema.Attachment{
		ID:          uuid.New(),
		URL:        expectedURL,
		Description: description,
		MaterialID:  &materialID,
	}

	suite.uploader.On("UploadFile", mock.AnythingOfType("string"), mock.AnythingOfType("*multipart.FileHeader")).Return(expectedURL, nil)


	suite.attachmentRepo.On("Create", ctx, mock.AnythingOfType("*schema.Attachment")).Return(nil)


	createdAttachment, err := suite.attachmentUseCase.CreateAttachment(ctx, fileHeader, description, materialID)


	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), createdAttachment.ID)
	assert.Equal(suite.T(), attachment.URL, createdAttachment.URL)
	suite.attachmentRepo.AssertExpectations(suite.T())
}

func (suite *AttachmentUseCaseTestSuite) TestCreateAssignmentAttachment_Success() {
	ctx := context.Background()
	fileHeader := &multipart.FileHeader{Filename: "assignment.pdf", Size: 1 * 1024 * 1024} // Assuming 1MB size
	description := "Assignment Attachment Description"
	assignmentID,_ := uuid.NewV7()

	expectedURL := "http://example.com/test.pdf"
	attachment := schema.Attachment{
		ID:          uuid.New(),
		URL:        expectedURL,
		Description: description,
		AssignmentID: &assignmentID,
	}

	suite.uploader.On("UploadFile", mock.AnythingOfType("string"), mock.AnythingOfType("*multipart.FileHeader")).Return(expectedURL, nil)



	suite.attachmentRepo.On("Create", ctx, mock.AnythingOfType("*schema.Attachment")).Return(nil)

	createdAttachment, err := suite.attachmentUseCase.CreateAssignmentAttachment(ctx, fileHeader, description, assignmentID)


	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), createdAttachment.ID)
	assert.Equal(suite.T(), attachment.URL, createdAttachment.URL)
	suite.attachmentRepo.AssertExpectations(suite.T())
}

func (suite *AttachmentUseCaseTestSuite) TestCreateSubmissionAttachment_Success() {
	ctx := context.Background()
	fileHeader := &multipart.FileHeader{Filename: "submission.pdf", Size: 1 * 1024 * 1024} 
	description := "Submission Attachment Description"

	expectedURL := "http://example.com/submission.pdf"
	attachment := schema.Attachment{
		ID:          uuid.New(),
		URL:         expectedURL,
		Description: description,
	}

	suite.uploader.On("UploadFile", mock.AnythingOfType("string"), mock.AnythingOfType("*multipart.FileHeader")).Return(expectedURL, nil)


	suite.attachmentRepo.On("Create", ctx, mock.AnythingOfType("*schema.Attachment")).Return(nil)

	createdAttachment, err := suite.attachmentUseCase.CreateSubmissionAttachment(ctx, fileHeader, description)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), createdAttachment.ID)
	assert.Equal(suite.T(), attachment.URL, createdAttachment.URL)
	suite.attachmentRepo.AssertExpectations(suite.T())
}

func (suite *AttachmentUseCaseTestSuite) TestGetAttachmentByID_Success() {
	ctx := context.Background()
	id := uuid.New()
	expectedAttachment := &schema.Attachment{
		ID:          id,
		URL:         "http://example.com/file.pdf",
		Description: "A test file",
	}

	suite.attachmentRepo.On("GetByID", ctx, id).Return(expectedAttachment, nil)

	attachment, err := suite.attachmentUseCase.GetAttachmentByID(ctx, id)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedAttachment, attachment)
	suite.attachmentRepo.AssertExpectations(suite.T())
}


func (suite *AttachmentUseCaseTestSuite) TestUpdateAttachment_Success() {
	ctx := context.Background()
	id := uuid.New()
	fileHeader := &multipart.FileHeader{Filename: "newfile.pdf", Size: 1 * 1024 * 1024}
	req := AttachmentUpdateRequest{
		File:        fileHeader,
		Description: "Updated description",
	}
	originalAttachment := &schema.Attachment{
		ID:          id,
		URL:         "http://example.com/original.pdf",
		Description: "Original description",
	}

	expectedURL := "http://example.com/newfile.pdf"
	
	updatedAttachment := &schema.Attachment{
		ID:          id,
		URL:         expectedURL,
		Description: req.Description,
	}

	suite.uploader.On("UploadFile", mock.AnythingOfType("string"), mock.AnythingOfType("*multipart.FileHeader")).Return(expectedURL, nil)

	suite.attachmentRepo.On("GetByID", ctx, id).Return(originalAttachment, nil)
	suite.attachmentRepo.On("Update", ctx, updatedAttachment).Return(nil)

	attachment, err := suite.attachmentUseCase.UpdateAttachment(ctx, id, req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedAttachment.Description, attachment.Description)
	assert.Equal(suite.T(), updatedAttachment.URL, attachment.URL)
	suite.attachmentRepo.AssertExpectations(suite.T())
}


func (suite *AttachmentUseCaseTestSuite) TestDeleteAttachment_Success() {
	ctx := context.Background()
	id := uuid.New()

	suite.attachmentRepo.On("Delete", ctx, id).Return(nil)

	err := suite.attachmentUseCase.DeleteAttachment(ctx, id)

	assert.NoError(suite.T(), err)
	suite.attachmentRepo.AssertExpectations(suite.T())
}

func TestAttachmentUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(AttachmentUseCaseTestSuite))
}
