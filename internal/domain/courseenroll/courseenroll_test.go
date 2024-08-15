package courseenroll

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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

type CourseEnrollUseCaseTestSuite struct {
	suite.Suite
	enrollRepo *MockEnrollRepository
	enrollUseCase *UseCase

}

func (suite *CourseEnrollUseCaseTestSuite) SetupTest() {
	
	suite.enrollRepo = new(MockEnrollRepository)
	suite.enrollUseCase = NewUseCase(suite.enrollRepo)
	
}

func (suite *CourseEnrollUseCaseTestSuite) TestEnrollStudent() {
    ctx := context.Background()
    userID := uuid.New()
    courseID := uuid.New()

    suite.enrollRepo.On("Create", ctx, mock.AnythingOfType("*schema.CourseEnroll")).Return(nil)

    err := suite.enrollUseCase.EnrollStudent(ctx, userID, courseID)

    assert.NoError(suite.T(), err)
    suite.enrollRepo.AssertExpectations(suite.T())
}

// TestGetEnrollmentsByCourse tests retrieval of users enrolled in a course.
func (suite *CourseEnrollUseCaseTestSuite) TestGetEnrollmentsByCourse() {
    ctx := context.Background()
    courseID := uuid.New()
    expectedUsers := []schema.User{{ID: uuid.New()}, {ID: uuid.New()}}

    suite.enrollRepo.On("GetUsersByCourseID", ctx, courseID).Return(expectedUsers, nil)

    users, err := suite.enrollUseCase.GetEnrollmentsByCourse(ctx, courseID)

    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), expectedUsers, users)
    suite.enrollRepo.AssertExpectations(suite.T())
}

// TestGetEnrollmentsByUser tests retrieval of courses a user is enrolled in.
func (suite *CourseEnrollUseCaseTestSuite) TestGetEnrollmentsByUser() {
    ctx := context.Background()
    userID := uuid.New()
    expectedCourses := []schema.Course{{ID: uuid.New()}, {ID: uuid.New()}}

    suite.enrollRepo.On("GetCoursesByUserID", ctx, userID).Return(expectedCourses, nil)

    courses, err := suite.enrollUseCase.GetEnrollmentsByUser(ctx, userID)

    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), expectedCourses, courses)
    suite.enrollRepo.AssertExpectations(suite.T())
}

// TestCheckEnrollment tests the enrollment status check.
func (suite *CourseEnrollUseCaseTestSuite) TestCheckEnrollment() {
    ctx := context.Background()
    userID := uuid.New()
    courseID := uuid.New()

    suite.enrollRepo.On("IsEnrolled", ctx, userID, courseID).Return(true, nil)

    enrolled, err := suite.enrollUseCase.CheckEnrollment(ctx, userID, courseID)

    assert.NoError(suite.T(), err)
    assert.True(suite.T(), enrolled)
    suite.enrollRepo.AssertExpectations(suite.T())
}

func TestCourseEnrollUseCaseTestSuite(t *testing.T) {
    suite.Run(t, new(CourseEnrollUseCaseTestSuite))
}