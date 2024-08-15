package forum

import (
	"context"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/courseenroll"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

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

type MockForumRepository struct {
	mock.Mock
}

func (m *MockForumRepository) CreateDiscussion(discussion *schema.ForumDiscussion) error {
	args := m.Called(discussion)
	return args.Error(0)
}

func (m *MockForumRepository) GetDiscussionByID(id uuid.UUID) (*schema.ForumDiscussion, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schema.ForumDiscussion), args.Error(1)
}

func (m *MockForumRepository) GetDiscussionsByCourseID(courseID uuid.UUID, page int, limit int) ([]*schema.ForumDiscussion, int64, error) {
	args := m.Called(courseID, page, limit)
	return args.Get(0).([]*schema.ForumDiscussion), args.Get(1).(int64), args.Error(2)
}

func (m *MockForumRepository) UpdateDiscussion(discussion *schema.ForumDiscussion) error {
	args := m.Called(discussion)
	return args.Error(0)
}

func (m *MockForumRepository) DeleteDiscussion(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockForumRepository) CreateReply(reply *schema.ForumReply) error {
	args := m.Called(reply)
	return args.Error(0)
}

func (m *MockForumRepository) GetReplyByID(id uuid.UUID) (*schema.ForumReply, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schema.ForumReply), args.Error(1)
}

func (m *MockForumRepository) GetRepliesByDiscussionID(discussionID uuid.UUID, page int, limit int) ([]*schema.ForumReply, int64, error) {
	args := m.Called(discussionID, page, limit)
	return args.Get(0).([]*schema.ForumReply), args.Get(1).(int64), args.Error(2)
}

func (m *MockForumRepository) UpdateReply(reply *schema.ForumReply) error {
	args := m.Called(reply)
	return args.Error(0)
}

func (m *MockForumRepository) DeleteReply(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

type ForumUseCaseTestSuite struct {
	suite.Suite
	forumRepo     *MockForumRepository
	courseRepo    *MockCourseRepository
	enrollRepo    *MockEnrollRepository
	enrollUseCase *courseenroll.UseCase
	forumUseCase  *UseCase
}

func (suite *ForumUseCaseTestSuite) SetupTest() {
	suite.forumRepo = new(MockForumRepository)
	suite.courseRepo = new(MockCourseRepository)
	suite.enrollRepo = new(MockEnrollRepository)
	suite.enrollUseCase = courseenroll.NewUseCase(suite.enrollRepo)
	suite.forumUseCase = NewUseCase(suite.forumRepo, suite.enrollUseCase, suite.courseRepo)
}

func (suite *ForumUseCaseTestSuite) TestCreateDiscussion_Success() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	ctx = context.WithValue(ctx, "user.role", "student")
	req := &CreateForumDiscussionRequest{
		Title:   "Test Discussion",
		Content: "This is a test discussion.",
	}

	suite.enrollRepo.On("IsEnrolled", ctx, userID, mock.Anything).Return(true, nil)
	suite.forumRepo.On("CreateDiscussion", mock.Anything).Return(nil)

	err := suite.forumUseCase.CreateDiscussion(ctx, req)

	assert.NoError(suite.T(), err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestCreateDiscussion_InvalidUserID() {
	ctx := context.WithValue(context.Background(), "user.id", "invalid-uuid")
	req := &CreateForumDiscussionRequest{
		Title:   "Test Discussion",
		Content: "This is a test discussion.",
	}

	err := suite.forumUseCase.CreateDiscussion(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrTokenInvalid, err)
}

func (suite *ForumUseCaseTestSuite) TestGetDiscussionByID_Success() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	ctx = context.WithValue(ctx, "user.role", "student")

	discussionID, _ := uuid.NewV7()
	discussion := &schema.ForumDiscussion{
		ID:      discussionID,
		Title:   "Test Discussion",
		Content: "This is a test discussion.",
	}

	suite.enrollRepo.On("IsEnrolled", ctx, userID, mock.Anything).Return(true, nil)
	suite.forumRepo.On("GetDiscussionByID", discussionID).Return(discussion, nil)

	res, err := suite.forumUseCase.GetDiscussionByID(ctx, discussionID.String())

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), res)
	assert.Equal(suite.T(), discussion, res)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestGetDiscussionByID_NotFound() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	ctx = context.WithValue(ctx, "user.role", "student")

	discussionID, _ := uuid.NewV7()

	suite.enrollRepo.On("IsEnrolled", ctx, userID, mock.Anything).Return(true, nil)
	suite.forumRepo.On("GetDiscussionByID", discussionID).Return(nil, gorm.ErrRecordNotFound)

	res, err := suite.forumUseCase.GetDiscussionByID(ctx, discussionID.String())

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), res)
	assert.Equal(suite.T(), ErrDiscussionNotFound, err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestUpdateDiscussion_Success() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())

	discussionID, _ := uuid.NewV7()
	req := &UpdateForumDiscussionRequest{
		ID:      discussionID.String(),
		Title:   "Updated Title",
		Content: "Updated content.",
	}

	discussion := &schema.ForumDiscussion{
		ID:      discussionID,
		UserID:  userID,
		Title:   "Test Discussion",
		Content: "This is a test discussion.",
	}

	suite.forumRepo.On("GetDiscussionByID", discussionID).Return(discussion, nil)
	suite.forumRepo.On("UpdateDiscussion", mock.Anything).Return(nil)

	err := suite.forumUseCase.UpdateDiscussion(ctx, req)

	assert.NoError(suite.T(), err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestUpdateDiscussion_NotFound() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())

	discussionID, _ := uuid.NewV7()
	req := &UpdateForumDiscussionRequest{
		ID:      discussionID.String(),
		Title:   "Updated Title",
		Content: "Updated content.",
	}

	suite.forumRepo.On("GetDiscussionByID", discussionID).Return(nil, gorm.ErrRecordNotFound)

	err := suite.forumUseCase.UpdateDiscussion(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrDiscussionNotFound, err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestUpdateDiscussion_NotYourResource() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())

	discussionID, _ := uuid.NewV7()
	req := &UpdateForumDiscussionRequest{
		ID:      discussionID.String(),
		Title:   "Updated Title",
		Content: "Updated content.",
	}

	discussion := &schema.ForumDiscussion{
		ID:      discussionID,
		UserID:  uuid.New(),
		Title:   "Test Discussion",
		Content: "This is a test discussion.",
	}

	suite.forumRepo.On("GetDiscussionByID", discussionID).Return(discussion, nil)

	err := suite.forumUseCase.UpdateDiscussion(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrNotYourResource, err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestDeleteDiscussion_Success() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	discussionID, _ := uuid.NewV7()

	discussion := &schema.ForumDiscussion{
		ID:      discussionID,
		UserID:  userID,
		Title:   "Test Discussion",
		Content: "This is a test discussion.",
	}

	suite.forumRepo.On("GetDiscussionByID", discussionID).Return(discussion, nil)
	suite.forumRepo.On("DeleteDiscussion", discussionID).Return(nil)

	err := suite.forumUseCase.DeleteDiscussion(ctx, discussionID.String())

	assert.NoError(suite.T(), err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestDeleteDiscussion_NotFound() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	discussionID, _ := uuid.NewV7()

	suite.forumRepo.On("GetDiscussionByID", discussionID).Return(nil, gorm.ErrRecordNotFound)

	err := suite.forumUseCase.DeleteDiscussion(ctx, discussionID.String())

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrDiscussionNotFound, err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestDeleteDiscussion_NotYourResource() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	discussionID, _ := uuid.NewV7()

	discussion := &schema.ForumDiscussion{
		ID:      discussionID,
		UserID:  uuid.New(),
		Title:   "Test Discussion",
		Content: "This is a test discussion.",
	}

	suite.forumRepo.On("GetDiscussionByID", discussionID).Return(discussion, nil)

	err := suite.forumUseCase.DeleteDiscussion(ctx, discussionID.String())

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrNotYourResource, err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestCreateReply_Success() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	ctx = context.WithValue(ctx, "user.role", "student")

	req := &CreateForumReplyRequest{
		DiscussionID: uuid.New(),
		Content:      "This is a reply",
	}

	suite.enrollRepo.On("IsEnrolled", ctx, userID, mock.Anything).Return(true, nil)
	suite.forumRepo.On("GetDiscussionByID", req.DiscussionID).Return(&schema.ForumDiscussion{}, nil)
	suite.forumRepo.On("CreateReply", mock.Anything).Return(nil)

	err := suite.forumUseCase.CreateReply(ctx, req)

	assert.NoError(suite.T(), err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestCreateReply_InvalidUserID() {
	ctx := context.WithValue(context.Background(), "user.id", "invalid-uuid")
	req := &CreateForumReplyRequest{
		DiscussionID: uuid.New(),
		Content:      "This is a reply",
	}

	suite.forumRepo.On("GetDiscussionByID", req.DiscussionID).Return(&schema.ForumDiscussion{}, nil)
	err := suite.forumUseCase.CreateReply(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrTokenInvalid, err)
}

func (suite *ForumUseCaseTestSuite) TestCreateReply_DiscussionNotFound() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	ctx = context.WithValue(ctx, "user.role", "student")

	req := &CreateForumReplyRequest{
		DiscussionID: uuid.New(),
		Content:      "This is a reply",
	}

	suite.enrollRepo.On("IsEnrolled", ctx, userID, mock.Anything).Return(true, nil)
	suite.forumRepo.On("GetDiscussionByID", req.DiscussionID).Return(nil, gorm.ErrRecordNotFound)

	err := suite.forumUseCase.CreateReply(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrDiscussionNotFound, err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestUpdateReply_Success() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	req := &UpdateForumReplyRequest{
		ID:      uuid.NewString(),
		Content: "Updated content",
	}

	reply := &schema.ForumReply{
		ID:     uuid.MustParse(req.ID),
		UserID: userID,
	}

	suite.forumRepo.On("GetReplyByID", reply.ID).Return(reply, nil)
	suite.forumRepo.On("UpdateReply", mock.Anything).Return(nil)

	err := suite.forumUseCase.UpdateReply(ctx, req)

	assert.NoError(suite.T(), err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestUpdateReply_NotFound() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	req := &UpdateForumReplyRequest{
		ID:      uuid.NewString(),
		Content: "Updated content",
	}

	suite.forumRepo.On("GetReplyByID", uuid.MustParse(req.ID)).Return(nil, gorm.ErrRecordNotFound)

	err := suite.forumUseCase.UpdateReply(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrReplyNotFound, err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestUpdateReply_NotYourResource() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	req := &UpdateForumReplyRequest{
		ID:      uuid.NewString(),
		Content: "Updated content",
	}

	reply := &schema.ForumReply{
		ID:     uuid.MustParse(req.ID),
		UserID: uuid.New(),
	}

	suite.forumRepo.On("GetReplyByID", reply.ID).Return(reply, nil)

	err := suite.forumUseCase.UpdateReply(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrNotYourResource, err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestDeleteReply_Success() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())

	replyID, _ := uuid.NewV7()

	suite.forumRepo.On("GetReplyByID", replyID).Return(&schema.ForumReply{ID: replyID, UserID: userID}, nil)
	suite.forumRepo.On("DeleteReply", replyID).Return(nil)

	err := suite.forumUseCase.DeleteReply(ctx, replyID.String())

	assert.NoError(suite.T(), err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestDeleteReply_NotFound() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())

	replyID, _ := uuid.NewV7()

	suite.forumRepo.On("GetReplyByID", replyID).Return(nil, gorm.ErrRecordNotFound)

	err := suite.forumUseCase.DeleteReply(ctx, replyID.String())

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrReplyNotFound, err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func (suite *ForumUseCaseTestSuite) TestDeleteReply_NotYourResource() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())

	replyID, _ := uuid.NewV7()

	suite.forumRepo.On("GetReplyByID", replyID).Return(&schema.ForumReply{ID: replyID, UserID: uuid.New()}, nil)

	err := suite.forumUseCase.DeleteReply(ctx, replyID.String())

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrNotYourResource, err)
	suite.forumRepo.AssertExpectations(suite.T())
}

func TestForumUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ForumUseCaseTestSuite))
}
