package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/fileutil"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"mime/multipart"
	"testing"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(user *schema.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) GetByID(id uuid.UUID) (*schema.User, error) {
	args := m.Called(id)
	user, ok := args.Get(0).(*schema.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockRepository) GetByEmail(email string) (*schema.User, error) {
	args := m.Called(email)
	user, ok := args.Get(0).(*schema.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockRepository) Update(user *schema.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) UpdateByEmail(email string, user *schema.User) error {
	args := m.Called(email, user)
	return args.Error(0)
}

type UseCaseTestSuite struct {
	suite.Suite
	repo    *MockRepository
	useCase *UseCase
}

func (suite *UseCaseTestSuite) SetupTest() {
	suite.repo = new(MockRepository)
	suite.useCase = NewUseCase(suite.repo)
}

func (suite *UseCaseTestSuite) TestGetMe_Success() {
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	user := &schema.User{ID: userID}

	suite.repo.On("GetByID", userID).Return(user, nil)

	result, err := suite.useCase.GetMe(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user, result)
}

func (suite *UseCaseTestSuite) TestGetMe_InvalidUUID() {
	ctx := context.WithValue(context.Background(), "user.id", "invalid-uuid")

	result, err := suite.useCase.GetMe(ctx)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), apierror.ErrTokenInvalid, err)
}

func (suite *UseCaseTestSuite) TestGetMe_InternalServerError() {
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())

	suite.repo.On("GetByID", userID).Return(nil, gorm.ErrRecordNotFound)

	result, err := suite.useCase.GetMe(ctx)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), apierror.ErrInternalServer, err)
}

func (suite *UseCaseTestSuite) TestGetByID_Success() {
	userID := uuid.New()
	req := &GetUserByIDRequest{ID: userID.String()}
	user := &schema.User{ID: userID}

	suite.repo.On("GetByID", userID).Return(user, nil)

	result, err := suite.useCase.GetByID(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), &GetUserResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		ImageURL: user.ImageURL,
		Role:     string(user.Role),
	}, result)
}

func (suite *UseCaseTestSuite) TestGetByID_InvalidUUID() {
	req := &GetUserByIDRequest{ID: "invalid-uuid"}

	result, err := suite.useCase.GetByID(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), apierror.ErrInvalidParamId, err)
}

func (suite *UseCaseTestSuite) TestGetByID_UserNotFound() {
	userID := uuid.New()
	req := &GetUserByIDRequest{ID: userID.String()}

	suite.repo.On("GetByID", userID).Return(nil, gorm.ErrRecordNotFound)

	result, err := suite.useCase.GetByID(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), ErrUserNotFound, err)
}

func (suite *UseCaseTestSuite) TestGetByID_InternalServerError() {
	userID := uuid.New()
	req := &GetUserByIDRequest{ID: userID.String()}

	suite.repo.On("GetByID", userID).Return(nil, gorm.ErrInvalidDB)

	result, err := suite.useCase.GetByID(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), apierror.ErrInternalServer, err)
}

func (suite *UseCaseTestSuite) TestUpdate_Success() {
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	req := &UpdateUserRequest{Name: "New Name"}
	user := &schema.User{ID: userID, Name: req.Name}

	suite.repo.On("Update", user).Return(nil)

	err := suite.useCase.Update(ctx, req)
	assert.NoError(suite.T(), err)
}

func (suite *UseCaseTestSuite) TestUpdate_InvalidUUID() {
	ctx := context.WithValue(context.Background(), "user.id", "invalid-uuid")
	req := &UpdateUserRequest{Name: "New Name"}

	err := suite.useCase.Update(ctx, req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrTokenInvalid, err)
}

func (suite *UseCaseTestSuite) TestUpdate_InvalidImageFile() {
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	req := &UpdateUserRequest{
		Name: "New Name",
		ImageFile: &multipart.FileHeader{
			Filename: "invalid_image.txt",
			Size:     1024,
		},
	}

	suite.repo.On("Update", mock.Anything).Return(apierror.ErrInvalidFileType)

	err := suite.useCase.Update(ctx, req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInvalidFileType, err)
}

func (suite *UseCaseTestSuite) TestUpdate_OversizedImageFile() {
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	req := &UpdateUserRequest{
		Name: "New Name",
		ImageFile: &multipart.FileHeader{
			Filename: "large_image.jpg",
			Size:     3 * fileutil.MegaByte,
		},
	}
	suite.repo.On("Update", mock.Anything).Return(apierror.ErrFileTooLarge)

	expectedErr := apierror.ErrFileTooLarge
	apierror.AddPayload(&expectedErr, map[string]string{
		"max_size":      "2 MB",
		"received_size": fileutil.ByteToAppropriateUnit(req.ImageFile.Size),
	})

	err := suite.useCase.Update(ctx, req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *UseCaseTestSuite) TestUpdate_InternalServerError() {
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	req := &UpdateUserRequest{Name: "New Name"}
	user := &schema.User{ID: userID, Name: req.Name}

	suite.repo.On("Update", user).Return(gorm.ErrInvalidDB)

	err := suite.useCase.Update(ctx, req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInternalServer, err)
}

func TestUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UseCaseTestSuite))
}
