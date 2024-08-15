package notification

import (
	"context"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) Create(notification *schema.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func (m *MockNotificationRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]*schema.Notification, int64, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
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

type NotificationUseCaseTestSuite struct {
	suite.Suite
	repo    *MockNotificationRepository
	useCase *UseCase
}

func (suite *NotificationUseCaseTestSuite) SetupTest() {
	suite.repo = new(MockNotificationRepository)
	suite.useCase = NewUseCase(suite.repo)
}

func (suite *NotificationUseCaseTestSuite) TestGetMy_Success() {
	ctx := context.WithValue(context.Background(), "user.id", uuid.New().String())
	req := &GetByUserIDRequest{Limit: 10, Page: 1}
	userID, _ := uuid.Parse(ctx.Value("user.id").(string))
	notifications := []*schema.Notification{{ID: uuid.New(), Title: "Test Notification"}}
	total := int64(1)

	suite.repo.On("GetByUserID", userID, req.Limit, 0).Return(notifications, total, nil)

	res, err := suite.useCase.GetMy(ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), res)
	resData := res.Data.([]*schema.Notification)
	assert.Equal(suite.T(), 1, len(resData))
	assert.Equal(suite.T(), total, int64(res.Pagination.TotalData))
	suite.repo.AssertExpectations(suite.T())
}

func (suite *NotificationUseCaseTestSuite) TestGetMy_InvalidUserID() {
	ctx := context.WithValue(context.Background(), "user.id", "invalid-uuid")
	req := &GetByUserIDRequest{Limit: 10, Page: 1}

	_, err := suite.useCase.GetMy(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrTokenInvalid, err)
}

func (suite *NotificationUseCaseTestSuite) TestGetMy_RepoError() {
	ctx := context.WithValue(context.Background(), "user.id", uuid.New().String())
	req := &GetByUserIDRequest{Limit: 10, Page: 1}
	userID, _ := uuid.Parse(ctx.Value("user.id").(string))

	suite.repo.On("GetByUserID", userID, req.Limit, 0).Return(nil, int64(0), gorm.ErrInvalidDB)

	_, err := suite.useCase.GetMy(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInternalServer, err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *NotificationUseCaseTestSuite) TestGetUnreadCount_Success() {
	ctx := context.WithValue(context.Background(), "user.id", uuid.New().String())
	userID, _ := uuid.Parse(ctx.Value("user.id").(string))
	count := int64(5)

	suite.repo.On("GetUnreadCount", userID).Return(count, nil)

	res, err := suite.useCase.GetUnreadCount(ctx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count, res)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *NotificationUseCaseTestSuite) TestGetUnreadCount_InvalidUserID() {
	ctx := context.WithValue(context.Background(), "user.id", "invalid-uuid")

	_, err := suite.useCase.GetUnreadCount(ctx)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrTokenInvalid, err)
}

func (suite *NotificationUseCaseTestSuite) TestGetUnreadCount_RepoError() {
	ctx := context.WithValue(context.Background(), "user.id", uuid.New().String())
	userID, _ := uuid.Parse(ctx.Value("user.id").(string))

	suite.repo.On("GetUnreadCount", userID).Return(int64(0), gorm.ErrInvalidDB)

	_, err := suite.useCase.GetUnreadCount(ctx)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInternalServer, err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *NotificationUseCaseTestSuite) TestUpdateRead_Success() {
	notificationID := uuid.New().String()
	id, _ := uuid.Parse(notificationID)

	suite.repo.On("UpdateRead", id).Return(nil)

	err := suite.useCase.UpdateRead(notificationID)

	assert.NoError(suite.T(), err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *NotificationUseCaseTestSuite) TestUpdateRead_InvalidNotificationID() {
	err := suite.useCase.UpdateRead("invalid-uuid")

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInvalidParamId, err)
}

func (suite *NotificationUseCaseTestSuite) TestUpdateRead_RepoError() {
	notificationID := uuid.New().String()
	id, _ := uuid.Parse(notificationID)

	suite.repo.On("UpdateRead", id).Return(apierror.ErrInternalServer)

	err := suite.useCase.UpdateRead(notificationID)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInternalServer, err)
	suite.repo.AssertExpectations(suite.T())
}

func TestNotificationUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationUseCaseTestSuite))
}
