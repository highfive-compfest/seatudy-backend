package review

import (
	"context"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/course"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/courseenroll"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type MockReviewRepository struct {
	mock.Mock
}

func (m *MockReviewRepository) Create(review *schema.Review, newCourseRating float32, newCourseRatingCount int64) error {
	args := m.Called(review, newCourseRating, newCourseRatingCount)
	return args.Error(0)
}

func (m *MockReviewRepository) GetByID(id uuid.UUID) (*schema.Review, error) {
	args := m.Called(id)
	var review *schema.Review
	if args.Get(0) != nil {
		review = args.Get(0).(*schema.Review)
	}
	return review, args.Error(1)
}

func (m *MockReviewRepository) Get(condition map[string]any, page int, limit int) ([]schema.Review, int64, error) {
	args := m.Called(condition, page, limit)
	return args.Get(0).([]schema.Review), args.Get(1).(int64), args.Error(2)
}

func (m *MockReviewRepository) Update(review *schema.Review, courseID uuid.UUID, newCourseRating float32) error {
	args := m.Called(review, courseID, newCourseRating)
	return args.Error(0)
}

func (m *MockReviewRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
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
	args := m.Called(ctx,courseID,userID)
	return args.Get(0).(float64),args.Error(1)
}

func (m *MockCourseRepository) SearchByTitle(ctx context.Context, title string, page, pageSize int) ([]schema.Course, int, error){
	args := m.Called(ctx, title, page, pageSize)
	return args.Get(0).([]schema.Course), args.Int(1), args.Error(2)
}

func (m *MockCourseRepository) DynamicFilterCourses(ctx context.Context, filterType, filterValue, sort string, page, limit int) ([]schema.Course, int, error){
	args := m.Called(ctx, filterType,filterValue,sort, page, limit)
	return args.Get(0).([]schema.Course), args.Int(1), args.Error(2)
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

type ReviewUseCaseTestSuite struct {
	suite.Suite
	reviewRepo    *MockReviewRepository
	courseRepo    *MockCourseRepository
	enrollRepo    *MockEnrollRepository
	reviewUseCase *UseCase
	enrollUseCase *courseenroll.UseCase
}

func (suite *ReviewUseCaseTestSuite) SetupTest() {
	suite.courseRepo = new(MockCourseRepository)
	suite.enrollRepo = new(MockEnrollRepository)
	suite.enrollUseCase = courseenroll.NewUseCase(suite.enrollRepo)
	suite.reviewRepo = new(MockReviewRepository)
	suite.reviewUseCase = NewUseCase(suite.reviewRepo, suite.courseRepo, suite.enrollUseCase)
}

func (suite *ReviewUseCaseTestSuite) TestCreateReview_Success() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	req := &CreateReviewRequest{
		CourseID: uuid.New(),
		Rating:   5,
		Feedback: "Great course!",
	}

	suite.enrollRepo.On("IsEnrolled", ctx, userID, req.CourseID).Return(true, nil)
	suite.courseRepo.On("GetRating", mock.Anything, req.CourseID).Return(float32(4.0), int64(10), nil)
	suite.reviewRepo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	res, err := suite.reviewUseCase.Create(ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), res)
	suite.courseRepo.AssertExpectations(suite.T())
	suite.reviewRepo.AssertExpectations(suite.T())
}

func (suite *ReviewUseCaseTestSuite) TestCreateReview_InvalidUserID() {
	ctx := context.WithValue(context.Background(), "user.id", "invalid-uuid")
	req := &CreateReviewRequest{
		CourseID: uuid.New(),
		Rating:   5,
		Feedback: "Great course!",
	}

	res, err := suite.reviewUseCase.Create(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), res)
	assert.Equal(suite.T(), apierror.ErrTokenInvalid, err)
}

func (suite *ReviewUseCaseTestSuite) TestCreateReview_NotEnrolled() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	req := &CreateReviewRequest{
		CourseID: uuid.New(),
		Rating:   5,
		Feedback: "Great course!",
	}

	suite.enrollRepo.On("IsEnrolled", ctx, userID, req.CourseID).Return(false, nil)

	res, err := suite.reviewUseCase.Create(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), res)
	assert.Equal(suite.T(), courseenroll.ErrNotEnrolled, err)
	suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *ReviewUseCaseTestSuite) TestCreateReview_CourseNotFound() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	req := &CreateReviewRequest{
		CourseID: uuid.New(),
		Rating:   5,
		Feedback: "Great course!",
	}

	suite.enrollRepo.On("IsEnrolled", ctx, userID, req.CourseID).Return(true, nil)
	suite.courseRepo.On("GetRating", mock.Anything, req.CourseID).Return(float32(0), int64(0), gorm.ErrRecordNotFound)

	res, err := suite.reviewUseCase.Create(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), res)
	assert.Equal(suite.T(), course.ErrCourseNotFound, err)
	suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *ReviewUseCaseTestSuite) TestCreateReview_AlreadyReviewed() {
	userID, _ := uuid.NewV7()
	ctx := context.WithValue(context.Background(), "user.id", userID.String())
	req := &CreateReviewRequest{
		CourseID: uuid.New(),
		Rating:   5,
		Feedback: "Great course!",
	}

	suite.enrollRepo.On("IsEnrolled", ctx, userID, req.CourseID).Return(true, nil)
	suite.courseRepo.On("GetRating", ctx, req.CourseID).Return(float32(4.0), int64(10), nil)

	pgErr := pgconn.PgError{Code: "23505"}

	suite.reviewRepo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(&pgErr)

	res, err := suite.reviewUseCase.Create(ctx, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), res)
	assert.Equal(suite.T(), ErrCourseAlreadyReviewed, err)
	suite.courseRepo.AssertExpectations(suite.T())
	suite.reviewRepo.AssertExpectations(suite.T())
}

func (suite *ReviewUseCaseTestSuite) TestGetReviews_Success() {
	ctx := context.Background()
	req := &GetReviewsRequest{
		CourseID: uuid.NewString(),
		Rating:   5,
		Page:     1,
		Limit:    10,
	}

	reviews := []schema.Review{
		{ID: uuid.New(), Rating: 5, Feedback: "Great course!"},
	}
	suite.reviewRepo.On("Get", mock.Anything, req.Page, req.Limit).Return(reviews, int64(1), nil)

	res, err := suite.reviewUseCase.Get(ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), res)
	assert.Equal(suite.T(), 1, len(res.Data.([]schema.Review)))
	suite.reviewRepo.AssertExpectations(suite.T())
}

func (suite *ReviewUseCaseTestSuite) TestUpdateReview_Success() {
	ctx := context.WithValue(context.Background(), "user.id", uuid.NewString())
	req := &UpdateReviewRequest{
		ID:       uuid.NewString(),
		Rating:   4,
		Feedback: "Updated feedback",
	}

	review := &schema.Review{
		ID:       uuid.MustParse(req.ID),
		UserID:   uuid.MustParse(ctx.Value("user.id").(string)),
		CourseID: uuid.New(),
		Rating:   5,
		Feedback: "Great course!",
	}

	suite.reviewRepo.On("GetByID", review.ID).Return(review, nil)
	suite.courseRepo.On("GetRating", mock.Anything, review.CourseID).Return(float32(4.0), int64(10), nil)
	suite.reviewRepo.On("Update", mock.Anything, review.CourseID, mock.Anything).Return(nil)

	err := suite.reviewUseCase.Update(ctx, req)

	assert.NoError(suite.T(), err)
	suite.reviewRepo.AssertExpectations(suite.T())
	suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *ReviewUseCaseTestSuite) TestUpdateReview_NotFound() {
	ctx := context.WithValue(context.Background(), "user.id", uuid.NewString())
	req := &UpdateReviewRequest{
		ID:       uuid.NewString(),
		Rating:   4,
		Feedback: "Updated feedback",
	}

	suite.reviewRepo.On("GetByID", uuid.MustParse(req.ID)).Return(nil, gorm.ErrRecordNotFound)

	err := suite.reviewUseCase.Update(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrReviewNotFound, err)
	suite.reviewRepo.AssertExpectations(suite.T())
}

func (suite *ReviewUseCaseTestSuite) TestUpdateReview_NotYourResource() {
	ctx := context.WithValue(context.Background(), "user.id", uuid.NewString())
	req := &UpdateReviewRequest{
		ID:       uuid.NewString(),
		Rating:   4,
		Feedback: "Updated feedback",
	}

	review := &schema.Review{
		ID:       uuid.MustParse(req.ID),
		UserID:   uuid.New(),
		CourseID: uuid.New(),
		Rating:   5,
		Feedback: "Great course!",
	}

	suite.reviewRepo.On("GetByID", review.ID).Return(review, nil)

	err := suite.reviewUseCase.Update(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrNotYourResource, err)
	suite.reviewRepo.AssertExpectations(suite.T())
}

func (suite *ReviewUseCaseTestSuite) TestDeleteReview_Success() {
	ctx := context.WithValue(context.Background(), "user.id", uuid.NewString())
	req := &DeleteReviewRequest{
		ID: uuid.NewString(),
	}

	review := &schema.Review{
		ID:       uuid.MustParse(req.ID),
		UserID:   uuid.MustParse(ctx.Value("user.id").(string)),
		CourseID: uuid.New(),
		Rating:   5,
		Feedback: "Great course!",
	}

	suite.reviewRepo.On("GetByID", review.ID).Return(review, nil)
	suite.reviewRepo.On("Delete", review.ID).Return(nil)

	err := suite.reviewUseCase.Delete(ctx, req)

	assert.NoError(suite.T(), err)
	suite.reviewRepo.AssertExpectations(suite.T())
}

func (suite *ReviewUseCaseTestSuite) TestDeleteReview_NotFound() {
	ctx := context.WithValue(context.Background(), "user.id", uuid.NewString())
	req := &DeleteReviewRequest{
		ID: uuid.NewString(),
	}

	suite.reviewRepo.On("GetByID", uuid.MustParse(req.ID)).Return(nil, gorm.ErrRecordNotFound)

	err := suite.reviewUseCase.Delete(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrReviewNotFound, err)
	suite.reviewRepo.AssertExpectations(suite.T())
}

func (suite *ReviewUseCaseTestSuite) TestDeleteReview_NotYourResource() {
	ctx := context.WithValue(context.Background(), "user.id", uuid.NewString())
	req := &DeleteReviewRequest{
		ID: uuid.NewString(),
	}

	review := &schema.Review{
		ID:       uuid.MustParse(req.ID),
		UserID:   uuid.New(),
		CourseID: uuid.New(),
		Rating:   5,
		Feedback: "Great course!",
	}

	suite.reviewRepo.On("GetByID", review.ID).Return(review, nil)

	err := suite.reviewUseCase.Delete(ctx, req)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrNotYourResource, err)
	suite.reviewRepo.AssertExpectations(suite.T())
}

func TestReviewUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ReviewUseCaseTestSuite))
}
