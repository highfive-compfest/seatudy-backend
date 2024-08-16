package material

import (
	"context"
	"mime/multipart"
	"testing"

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

// Create mocks the Create method
func (m *MockRepository) Create(ctx context.Context, mat *schema.Material) error {
	args := m.Called(ctx, mat)
	return args.Error(0)
}

// GetByID mocks the GetByID method
func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*schema.Material, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*schema.Material), args.Error(1)
}

// GetAll mocks the GetAll method
func (m *MockRepository) GetAll(ctx context.Context) ([]*schema.Material, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*schema.Material), args.Error(1)
}

// Update mocks the Update method
func (m *MockRepository) Update(ctx context.Context, mat *schema.Material) error {
	args := m.Called(ctx, mat)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockAttachmentRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockAttachmentRepository) Create(ctx context.Context, att *schema.Attachment) error {
	args := m.Called(ctx, att)
	return args.Error(0)
}

// Update mocks the Update method
func (m *MockAttachmentRepository) Update(ctx context.Context, att *schema.Attachment) error {
	args := m.Called(ctx, att)
	return args.Error(0)
}

// GetByID mocks the GetByID method
func (m *MockAttachmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*schema.Attachment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*schema.Attachment), args.Error(1)
}

// Delete mocks the Delete method
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

type MaterialUseCaseTestSuite struct {
	suite.Suite
	attachmentRepo    *MockAttachmentRepository
	uploader          *MockFileUploader
	attachmentUseCase *attachment.UseCase
	materialRepo      *MockRepository
	materialUseCase   *UseCase
}

func (suite *MaterialUseCaseTestSuite) SetupTest() {
	suite.attachmentRepo = new(MockAttachmentRepository)
	suite.uploader = new(MockFileUploader)
	suite.materialRepo = new(MockRepository)
	suite.attachmentUseCase = attachment.NewUseCase(suite.attachmentRepo, suite.uploader)
	suite.materialUseCase = NewUseCase(suite.materialRepo,suite.attachmentUseCase)
}

func (suite *MaterialUseCaseTestSuite) TestCreateMaterial_Success() {
	ctx := context.Background()
	req := CreateMaterialRequest{
		CourseID:    uuid.New().String(),
		Title:       "Introduction to Chemistry",
		Description: "A basic introduction to chemistry principles.",
	}

	suite.materialRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Call the function under test
	err := suite.materialUseCase.CreateMaterial(ctx, req)

	// Assertions
	assert.NoError(suite.T(), err)
	suite.materialRepo.AssertExpectations(suite.T())
}


func (suite *MaterialUseCaseTestSuite) TestUpdateMaterial_Success() {
	ctx := context.Background()
	materialID := uuid.New()
	existingMaterial := &schema.Material{
		ID:          materialID,
		Title:       "Old Title",
		Description: "Old Description",
	}
	updatedTitle := "Updated Title"
	updatedDescription := "Updated Description"

	req := UpdateMaterialRequest{
		Title:       &updatedTitle,
		Description: &updatedDescription,
	}

	// Mocking the repository to return the existing material and handle the update
	suite.materialRepo.On("GetByID", ctx, materialID).Return(existingMaterial, nil)
	suite.materialRepo.On("Update", ctx, mock.AnythingOfType("*schema.Material")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*schema.Material)
		assert.Equal(suite.T(), updatedTitle, arg.Title)
		assert.Equal(suite.T(), updatedDescription, arg.Description)
	})

	// Call the function under test
	err := suite.materialUseCase.UpdateMaterial(ctx, req, materialID)

	// Assertions
	assert.NoError(suite.T(), err)
	suite.materialRepo.AssertExpectations(suite.T())
}

func (suite *MaterialUseCaseTestSuite) TestGetMaterialByID_Success() {
	ctx := context.Background()
	materialID := uuid.New()
	expectedMaterial := &schema.Material{
		ID:          materialID,
		Title:       "Chemistry Basics",
		Description: "A foundational course in chemistry.",
	}

	suite.materialRepo.On("GetByID", ctx, materialID).Return(expectedMaterial, nil)

	// Call the function under test
	material, err := suite.materialUseCase.GetMaterialByID(ctx, materialID)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedMaterial, material)
	suite.materialRepo.AssertExpectations(suite.T())
}

// TestGetAllMaterials_Success tests the successful retrieval of all materials
func (suite *MaterialUseCaseTestSuite) TestGetAllMaterials_Success() {
	ctx := context.Background()
	expectedMaterials := []*schema.Material{
		{ID: uuid.New(), Title: "Chemistry Basics", Description: "A foundational course in chemistry."},
		{ID: uuid.New(), Title: "Physics Introduction", Description: "Basic concepts of physics."},
	}

	suite.materialRepo.On("GetAll", ctx).Return(expectedMaterials, nil)

	// Call the function under test
	materials, err := suite.materialUseCase.GetAllMaterials(ctx)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedMaterials, materials)
	suite.materialRepo.AssertExpectations(suite.T())
}

func (suite *MaterialUseCaseTestSuite) TestDeleteMaterial_Success() {
	ctx := context.Background()
	materialID := uuid.New()

	suite.materialRepo.On("Delete", ctx, materialID).Return(nil)

	err := suite.materialUseCase.DeleteMaterial(ctx, materialID)

	assert.NoError(suite.T(), err)
	suite.materialRepo.AssertExpectations(suite.T())
}

func (suite *MaterialUseCaseTestSuite) TestAddAttachment_Success() {
	ctx := context.Background()
	materialID := uuid.New()
	fileHeader := &multipart.FileHeader{Filename: "test_attachment.pdf", Size: 1024}
	req := AttachmentInput{
		File:        fileHeader,
		Description: "Test attachment description",
	}

	expectedAttachment := &schema.Attachment{
		ID:          uuid.New(),
		Description: req.Description,
		URL:         "http://example.com/test_attachment.pdf",
	}

	existingMaterial := &schema.Material{
		ID:          materialID,
		Title:       "Existing Material",
		Attachments: []schema.Attachment{},
	}


	suite.materialRepo.On("GetByID", ctx, materialID).Return(existingMaterial, nil)
	suite.uploader.On("UploadFile", mock.AnythingOfType("string"), mock.AnythingOfType("*multipart.FileHeader")).Return(expectedAttachment.URL, nil)
	suite.attachmentRepo.On("Create", ctx, mock.AnythingOfType("*schema.Attachment")).Return(nil)
	suite.materialRepo.On("Update", ctx, existingMaterial).Return(nil)

	err := suite.materialUseCase.AddAttachment(ctx, materialID, req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(existingMaterial.Attachments))
	assert.Equal(suite.T(), expectedAttachment.URL, existingMaterial.Attachments[0].URL) 
	suite.materialRepo.AssertExpectations(suite.T())
	suite.attachmentRepo.AssertExpectations(suite.T())
}


func TestMaterialUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(MaterialUseCaseTestSuite))
}