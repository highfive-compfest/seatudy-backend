package course

import (
	"context"
	"mime/multipart"
	"os"
	"testing"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"

	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/config"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/courseenroll"
	"github.com/highfive-compfest/seatudy-backend/internal/fileutil"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

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

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) Create(tx *gorm.DB, wallet *schema.Wallet) error {
	args := m.Called(tx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) CreateMidtransTransaction(tx *gorm.DB, transaction *schema.MidtransTransaction) error {
	args := m.Called(tx, transaction)
	return args.Error(0)
}

func (m *MockWalletRepository) GetByUserID(tx *gorm.DB, userID uuid.UUID) (*schema.Wallet, error) {
	args := m.Called(tx, userID)
	var wallet *schema.Wallet
	if args.Get(0) != nil {
		wallet = args.Get(0).(*schema.Wallet)
	}
	return wallet, args.Error(1)
}

func (m *MockWalletRepository) GetMidtransTransactionByID(tx *gorm.DB, transactionID uuid.UUID) (*schema.MidtransTransaction, error) {
	args := m.Called(tx, transactionID)
	return args.Get(0).(*schema.MidtransTransaction), args.Error(1)
}

func (m *MockWalletRepository) GetMidtransTransactionsByWalletID(tx *gorm.DB, walletID uuid.UUID, isCredit bool, page, limit int) ([]*schema.MidtransTransaction, int64, error) {
	args := m.Called(tx, walletID, isCredit, page, limit)
	var transactions []*schema.MidtransTransaction

	if args.Get(0) != nil {
		transactions = args.Get(0).([]*schema.MidtransTransaction)
	}

	return transactions, args.Get(1).(int64), args.Error(2)
}

func (m *MockWalletRepository) UpdateMidtransTransaction(tx *gorm.DB, transaction *schema.MidtransTransaction) error {
	args := m.Called(tx, transaction)
	return args.Error(0)
}

func (m *MockWalletRepository) TopUpSuccess(transactionID uuid.UUID) error {
	args := m.Called(transactionID)
	return args.Error(0)
}

func (m *MockWalletRepository) TransferByUserID(tx *gorm.DB, fromUserID, toUserID uuid.UUID, amount int64) error {
	args := m.Called(tx, fromUserID, toUserID, amount)
	return args.Error(0)
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

type MockFileUploader struct {
	mock.Mock
}

func (m *MockFileUploader) UploadFile(key string, fileHeader *multipart.FileHeader) (string, error) {
	args := m.Called(key, fileHeader)
	return args.String(0), args.Error(1)
}


type CourseUseCaseTestSuite struct {
	suite.Suite
	walletRepo    *MockWalletRepository
	courseRepo    *MockCourseRepository
	enrollRepo    *MockEnrollRepository
	userRepo *MockUserRepository
	mailer *MockMailer
	enrollUseCase *courseenroll.UseCase
	notificationRepo *MockNotificationRepository
	courseUseCase *UseCase
    uploader *MockFileUploader
}

func (suite *CourseUseCaseTestSuite) SetupTest() {
	suite.courseRepo = new(MockCourseRepository)
	suite.enrollRepo = new(MockEnrollRepository)
	suite.walletRepo = new(MockWalletRepository)
	suite.userRepo =  new(MockUserRepository)
    suite.uploader = new(MockFileUploader)
	suite.mailer = new(MockMailer)
	suite.notificationRepo = new(MockNotificationRepository)
	suite.enrollUseCase = courseenroll.NewUseCase(suite.enrollRepo)
	suite.courseUseCase = NewUseCase(suite.courseRepo,suite.walletRepo,*suite.enrollUseCase, suite.userRepo,suite.notificationRepo,suite.mailer, suite.uploader)

}

func (suite *CourseUseCaseTestSuite) TestCreateCourse_Success() {
    // Mock input
    ctx := context.Background()
    req := CreateCourseRequest{
        Title:       "Introduction to Go",
        Description: "Learn the basics of Go.",
        Price:       10000,
        Difficulty:  "beginner",
        Category:    "Programming Languages",
    }

    instructorID,_ := uuid.NewV7()

	
    suite.courseRepo.On("Create", ctx, mock.Anything).Return(nil)


    err := suite.courseUseCase.Create(ctx, req, nil, nil, instructorID.String())


	assert.NoError(suite.T(), err)
	suite.courseRepo.AssertExpectations(suite.T())

}

func (suite *CourseUseCaseTestSuite) TestCreateCourse_InvalidImageFileType() {
    // Mock input
    ctx := context.Background()
    req := CreateCourseRequest{
        Title:       "Introduction to Go",
        Description: "Learn the basics of Go.",
        Price:       10000,
        Difficulty:  "beginner",
        Category:    "Programming Languages",
    }
    imageFile := &multipart.FileHeader{Filename: "image.pdf", Size: 1 * fileutil.MegaByte} // Unsupported file type for this test
    instructorID, _ := uuid.NewV7()

    // This should return an error due to invalid image file type
	
    err := suite.courseUseCase.Create(ctx, req, imageFile, nil, instructorID.String())

	
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(),  apierror.ErrInternalServer, err)
	suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestCreateCourse_InvalidSyllabusFileType() {
    // Mock input
    ctx := context.Background()
    req := CreateCourseRequest{
        Title:       "Advanced Node.js",
        Description: "Deep dive into Node.js.",
        Price:       12000,
        Difficulty:  "advanced",
        Category:    "Programming Languages",
    }

    syllabusFile := &multipart.FileHeader{Filename: "syllabus.jpg", Size: 1 * fileutil.MegaByte} // Unsupported file type
    instructorID, _ := uuid.NewV7()


    err := suite.courseUseCase.Create(ctx, req, nil, syllabusFile, instructorID.String())

    assert.Error(suite.T(), err)
	assert.Equal(suite.T(),  apierror.ErrInternalServer, err)
	suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestCreateCourse_InvalidInstructorID() {
    // Mock input
    ctx := context.Background()
    req := CreateCourseRequest{
        Title:       "Python for Data Science",
        Description: "Explore Python in Data Science.",
        Price:       15000,
        Difficulty:  "intermediate",
        Category:    "Data Science & Analytics",
    }

    invalidId := "invalid"

    // This should return an error due to invalid instructor ID
    err := suite.courseUseCase.Create(ctx, req, nil,nil, invalidId)


    assert.Error(suite.T(), err)
	assert.Equal(suite.T(),  ErrUnauthorizedAccess , err)
	suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestCreateCourse_FileTooLarge() {
    // Mock input
    ctx := context.Background()
    req := CreateCourseRequest{
        Title:       "Python for Data Science",
        Description: "Explore Python in Data Science.",
        Price:       15000,
        Difficulty:  "intermediate",
        Category:    "Data Science & Analytics",
    }
    imageFile := &multipart.FileHeader{Filename: "image.jpg", Size: 3 * fileutil.MegaByte}
    syllabusFile := &multipart.FileHeader{Filename: "syllabus.pdf", Size: 1 * fileutil.MegaByte}
    instructorID, _ := uuid.NewV7() 

    // This should return an error due to invalid instructor ID
    err := suite.courseUseCase.Create(ctx, req, imageFile, syllabusFile, instructorID.String())

	err2 := apierror.ErrFileTooLarge
	apierror.AddPayload(&err2, map[string]string{
		"max_size":      "2 MB",
		"received_size": fileutil.ByteToAppropriateUnit(imageFile.Size),
	})
    assert.Error(suite.T(), err)
	assert.Equal(suite.T(),  err2, err)
	suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestGetAll_Success() {
    ctx := context.Background()
    page, pageSize := 1, 10
	courseID, _ := uuid.NewV7()
	courseID1, _ := uuid.NewV7()
    mockCourses := []schema.Course{{ID: courseID, Title: "Go Programming"}, {ID: courseID1, Title: "Python Programming"}}
    mockTotal := 2

    suite.courseRepo.On("GetAll", ctx, page, pageSize).Return(mockCourses, mockTotal, nil)

    response, err := suite.courseUseCase.GetAll(ctx, page, pageSize)

    assert.NoError(suite.T(), err)
    assert.Len(suite.T(), response.Courses, 2)
    assert.Equal(suite.T(), mockTotal, response.Pagination.TotalData)
    suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestGetCourseByPopularity_Success() {
    ctx := context.Background()
    page, pageSize := 1, 10
	courseID, _ := uuid.NewV7()
	courseID1, _ := uuid.NewV7()
    mockCourses := []schema.Course{{ID: courseID, Title: "Node.js Best Practices"}, {ID: courseID1, Title: "Advanced React"}}
    mockTotal := 2

    suite.courseRepo.On("FindByPopularity", ctx, page, pageSize).Return(mockCourses, mockTotal, nil)

    response, err := suite.courseUseCase.GetCourseByPopularity(ctx, page, pageSize)

    assert.NoError(suite.T(), err)
    assert.Len(suite.T(), response.Courses, 2)
    assert.Equal(suite.T(), mockTotal, response.Pagination.TotalData)
    suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestGetByInstructorID_Success() {
    ctx := context.Background()
    page, pageSize := 1, 10
    instructorID,_ := uuid.NewV7()
	courseID, _ := uuid.NewV7()
    mockCourses := []schema.Course{{ID: courseID, Title: "Intro to Docker"}}
    mockTotal := 1

    suite.courseRepo.On("FindByInstructorID", ctx, instructorID, page, pageSize).Return(mockCourses, mockTotal, nil)

    response, err := suite.courseUseCase.GetByInstructorID(ctx, instructorID, page, pageSize)

    assert.NoError(suite.T(), err)
    assert.Len(suite.T(), response.Courses, 1)
    suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestGetByID_Success() {
    ctx := context.Background()
    courseID := uuid.New()
    mockCourse := schema.Course{ID: courseID, Title: "Microservices with Go"}

    suite.courseRepo.On("GetByID", ctx, courseID).Return(mockCourse, nil)

    course, err := suite.courseUseCase.GetByID(ctx, courseID)

    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), mockCourse.Title, course.Title)
    suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestUpdate_Success() {
    ctx := context.Background()
    id := uuid.New()
    mockCourse := schema.Course{ID: id, Title: "Original Title"}
    req := UpdateCourseRequest{Title: new(string)}
    *req.Title = "Updated Title"

    suite.courseRepo.On("GetByID", ctx, id).Return(mockCourse, nil)
    suite.courseRepo.On("Update", ctx, mock.Anything).Return(nil)

    updatedCourse, err := suite.courseUseCase.Update(ctx, req, id, nil, nil)

    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), "Updated Title", updatedCourse.Title)
    suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestUpdate_InvalidFileType() {
    ctx := context.Background()
    id := uuid.New()
    mockCourse := schema.Course{ID: id, Title: "Original Title"}
    req := UpdateCourseRequest{}
    imageFile := &multipart.FileHeader{Filename: "image.exe", Size: 1 * fileutil.MegaByte}

    suite.courseRepo.On("GetByID", ctx, id).Return(mockCourse, nil)
    suite.courseRepo.On("Update", ctx, mock.Anything).Return(nil) // This should not be called

    _, err := suite.courseUseCase.Update(ctx, req, id, imageFile, nil)

    assert.Error(suite.T(), err)
    assert.IsType(suite.T(), apierror.ErrInternalServer, err)
}

func (suite *CourseUseCaseTestSuite) TestUpdate_InvalidCourseID() {
    ctx := context.Background()
    id,_ := uuid.NewV7() // Suppose this ID does not exist in the database
    req := UpdateCourseRequest{}

    suite.courseRepo.On("GetByID", ctx, id).Return(schema.Course{}, ErrCourseNotFound)

    _, err := suite.courseUseCase.Update(ctx, req, id, nil, nil)

    assert.Error(suite.T(), err)
    assert.IsType(suite.T(), ErrCourseNotFound, err)
}

func (suite *CourseUseCaseTestSuite) TestSearchCoursesByTitle_Success() {
    ctx := context.Background()
    title := "Go Programming"
    page, pageSize := 1, 10
    mockCourses := []schema.Course{
        {Title: "Go Programming 101", ID: uuid.New()},
        {Title: "Advanced Go Programming", ID: uuid.New()},
    }
    total := len(mockCourses)

    suite.courseRepo.On("SearchByTitle", ctx, title, page, pageSize).Return(mockCourses, total, nil)

    response, err := suite.courseUseCase.SearchCoursesByTitle(ctx, title, page, pageSize)

    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), total, response.Pagination.TotalData)
    assert.Len(suite.T(), response.Courses, total)
    suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestBuyCourse_Success() {
    ctx := context.Background()
	ctx = context.WithValue(ctx, "user.name", "John Doe")
    ctx = context.WithValue(ctx, "user.email", "john.doe@example.com")
    _ = os.Setenv("ENV", "test")
	_ = os.Setenv("SMTP_EMAIL", "test")
	config.LoadEnv()
    courseId, _ := uuid.NewV7()
    studentId, _ := uuid.NewV7()
    instructorId, _ := uuid.NewV7()

    mockCourse := schema.Course{ID: courseId, InstructorID: instructorId, Price: 10000}
	mockInstructor := schema.User{ID: instructorId, Name: "Instructor Name", Email: "instructor@example.com"}

    suite.courseRepo.On("GetByID", ctx, courseId).Return(mockCourse, nil)
	suite.userRepo.On("GetByID", instructorId).Return(&mockInstructor, nil)

    suite.enrollRepo.On("IsEnrolled", ctx, studentId, courseId).Return(false, nil)
    suite.mailer.On("DialAndSend", mock.Anything).Return(nil)
 
    suite.walletRepo.On("TransferByUserID", mock.AnythingOfType("*gorm.DB"), studentId, instructorId, int64(10000)).Return(nil)


    suite.enrollRepo.On("Create", ctx, mock.AnythingOfType("*schema.CourseEnroll")).Return(nil)
	suite.notificationRepo.On("Create", mock.AnythingOfType("*schema.Notification")).Return(nil)

    // Executing the method under test
    err := suite.courseUseCase.BuyCourse(ctx, courseId, studentId.String())

    // Assertions to check that no error occurred and all expectations were met
    assert.NoError(suite.T(), err)
    suite.courseRepo.AssertExpectations(suite.T())
    suite.enrollRepo.AssertExpectations(suite.T())
    suite.walletRepo.AssertExpectations(suite.T())

}


func (suite *CourseUseCaseTestSuite) TestBuyCourse_CourseNotFound() {
    ctx := context.Background()
    courseId,_ := uuid.NewV7()
    studentId,_ := uuid.NewV7()

    suite.courseRepo.On("GetByID", ctx, courseId).Return(schema.Course{}, gorm.ErrRecordNotFound)

    err := suite.courseUseCase.BuyCourse(ctx, courseId, studentId.String())

    assert.Error(suite.T(), err)
    assert.Equal(suite.T(), ErrCourseNotFound, err)
}

func (suite *CourseUseCaseTestSuite) TestBuyCourse_AlreadyEnrolled() {
    ctx := context.Background()
    courseId,_ := uuid.NewV7()
    studentId,_ := uuid.NewV7()

    mockCourse := schema.Course{ID: courseId}
    suite.courseRepo.On("GetByID", ctx, courseId).Return(mockCourse, nil)
    suite.enrollRepo.On("IsEnrolled", ctx, studentId, courseId).Return(true, nil)

    err := suite.courseUseCase.BuyCourse(ctx, courseId, studentId.String())

    assert.Error(suite.T(), err)
    assert.Equal(suite.T(), ErrAlreadyEnrolled, err)
}

func (suite *CourseUseCaseTestSuite) TestBuyCourse_TransactionError() {
    ctx := context.Background()
    courseId,_ := uuid.NewV7()
    studentId,_ := uuid.NewV7()
    instructorId := uuid.New()

    mockCourse := schema.Course{ID: courseId, InstructorID: instructorId, Price: 10000}
    suite.courseRepo.On("GetByID", ctx, courseId).Return(mockCourse, nil)
    suite.enrollRepo.On("IsEnrolled", ctx, studentId, courseId).Return(false, nil)
    suite.walletRepo.On("TransferByUserID", mock.Anything, studentId, instructorId, int64(10000)).Return(apierror.ErrInternalServer)

    err := suite.courseUseCase.BuyCourse(ctx, courseId, studentId.String())

    assert.Error(suite.T(), err)
    assert.Equal(suite.T(), apierror.ErrInternalServer, err)
}

func (suite *CourseUseCaseTestSuite) TestGetEnrollmentsByCourse_Success() {
    ctx := context.Background()
    courseId := uuid.New()

    // Mocking enrolled users
    mockUsers := []schema.User{
        {ID: uuid.New(), Name: "John Doe"},
        {ID: uuid.New(), Name: "Jane Doe"},
    }

    suite.enrollRepo.On("GetUsersByCourseID", ctx, courseId).Return(mockUsers, nil)

    users, err := suite.courseUseCase.GetEnrollmentsByCourse(ctx, courseId)

    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), mockUsers, users)
    suite.enrollRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestGetEnrollmentsByCourse_Failure() {
    ctx := context.Background()
    courseId := uuid.New()

    suite.enrollRepo.On("GetUsersByCourseID", ctx, courseId).Return([]schema.User{}, apierror.ErrInternalServer)

    users, err := suite.courseUseCase.GetEnrollmentsByCourse(ctx, courseId)

    assert.Error(suite.T(), err)
    assert.Nil(suite.T(), users)
    suite.enrollRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestGetEnrollmentsByUser_Success() {
    ctx := context.Background()
    studentId,_ := uuid.NewV7()



    // Mocking enrolled courses and user progress
    mockCourses := []schema.Course{
        {ID: uuid.New(), Title: "Course 1"},
        {ID: uuid.New(), Title: "Course 2"},
    }
    mockProgress := []float64{75.0, 90.0}

    suite.enrollRepo.On("GetCoursesByUserID", ctx, studentId).Return(mockCourses, nil)
    for i, course := range mockCourses {
		suite.courseRepo.On("GetByID", ctx, course.ID).Return(course, nil)
        suite.courseRepo.On("GetUserCourseProgress", ctx, course.ID, studentId).Return(mockProgress[i], nil)
    }

    courses, err := suite.courseUseCase.GetEnrollmentsByUser(ctx, studentId.String())

    assert.NoError(suite.T(), err)
    assert.Len(suite.T(), courses, len(mockCourses))
    for i, resp := range courses {
        assert.Equal(suite.T(), mockCourses[i].Title, resp.Course.Title)
        assert.Equal(suite.T(), mockProgress[i], resp.Progress)
    }
    suite.enrollRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestGetEnrollmentsByUser_Failure() {
    ctx := context.Background()
    studentId, _ := uuid.NewV7()

    // Return an empty slice and an error
    suite.enrollRepo.On("GetCoursesByUserID", ctx, studentId).Return([]schema.Course{}, apierror.ErrInternalServer)

    courses, err := suite.courseUseCase.GetEnrollmentsByUser(ctx, studentId.String())

    assert.Error(suite.T(), err)
    assert.Nil(suite.T(), courses) // This may need to be adjusted based on actual method implementation
    suite.enrollRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestGetUserCourseProgress_Success() {
    ctx := context.Background()
    courseId := uuid.New()
    userId := uuid.New().String()
    expectedProgress := 85.0

    mockCourse := schema.Course{ID: courseId}
    suite.courseRepo.On("GetByID", ctx, courseId).Return(mockCourse, nil)
    suite.courseRepo.On("GetUserCourseProgress", ctx, courseId, mock.Anything).Return(expectedProgress, nil)

    progress, err := suite.courseUseCase.GetUserCourseProgress(ctx, courseId, userId)

    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), expectedProgress, progress)
    suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestGetUserCourseProgress_CourseNotFound() {
    ctx := context.Background()
    courseId := uuid.New()
    userId := uuid.New().String()

    suite.courseRepo.On("GetByID", ctx, courseId).Return(schema.Course{}, ErrCourseNotFound)

    progress, err := suite.courseUseCase.GetUserCourseProgress(ctx, courseId, userId)

    assert.Error(suite.T(), err)
    assert.Equal(suite.T(), ErrCourseNotFound, err)
    assert.Equal(suite.T(), 0.0, progress)
    suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestGetUserCourseProgress_InvalidUserID() {
    ctx := context.Background()
    courseId := uuid.New()
    invalidUserId := "invalid-uuid"

    mockCourse := schema.Course{ID: courseId}
    suite.courseRepo.On("GetByID", ctx, courseId).Return(mockCourse, nil)

    progress, err := suite.courseUseCase.GetUserCourseProgress(ctx, courseId, invalidUserId)

    assert.Error(suite.T(), err)
    assert.Equal(suite.T(), apierror.ErrInternalServer, err)
    assert.Equal(suite.T(), 0.0, progress)
    suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestGetUserCourseProgress_FetchProgressError() {
    ctx := context.Background()
    courseId := uuid.New()
    userId := uuid.New().String()

    mockCourse := schema.Course{ID: courseId}
    suite.courseRepo.On("GetByID", ctx, courseId).Return(mockCourse, nil)
    suite.courseRepo.On("GetUserCourseProgress", ctx, courseId, mock.Anything).Return(0.0, apierror.ErrInternalServer)

    progress, err := suite.courseUseCase.GetUserCourseProgress(ctx, courseId, userId)

    assert.Error(suite.T(), err)
    assert.Equal(suite.T(), apierror.ErrInternalServer, err)
    assert.Equal(suite.T(), 0.0, progress)
    suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestFilterCourses_ByCategory() {
    ctx := context.Background()
    category := "Programming Languages"
    req := FilterCoursesRequest{
        Category: &category,
        Page:     1,
        Limit:    10,
    }
    mockCourses := []schema.Course{{ID: uuid.New(), Title: "Intro to Go"}, {ID: uuid.New(), Title: "Advanced Python"}}
    total := 2

    suite.courseRepo.On("DynamicFilterCourses", ctx, "category", "Programming Languages", "", 1, 10).Return(mockCourses, total, nil)

    response, err := suite.courseUseCase.FilterCourses(ctx, req)

    assert.NoError(suite.T(), err)
    assert.Len(suite.T(), response.Courses, 2)
    assert.Equal(suite.T(), total, response.Pagination.TotalData)
    suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestFilterCourses_ByDifficulty() {
    ctx := context.Background()
    difficulty := "beginner"
    req := FilterCoursesRequest{
        Difficulty: &difficulty,
        Page:       1,
        Limit:      10,
    }

    mockCourses := []schema.Course{{ID: uuid.New(), Title: "Basic HTML"}, {ID: uuid.New(), Title: "CSS Fundamentals"}}
    total := 2

    suite.courseRepo.On("DynamicFilterCourses", ctx, "difficulty", "beginner", "", 1, 10).Return(mockCourses, total, nil)

    response, err := suite.courseUseCase.FilterCourses(ctx, req)

    assert.NoError(suite.T(), err)
    assert.Len(suite.T(), response.Courses, 2)
    suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestFilterCourses_ByRating() {
    ctx := context.Background()
    rating := float32(4.0)
    req := FilterCoursesRequest{
        Rating: &rating,
        Page:   1,
        Limit:  10,
    }

    mockCourses := []schema.Course{{ID: uuid.New(), Title: "React for Beginners"}, {ID: uuid.New(), Title: "Vue Advanced Techniques"}}
    total := 2

    suite.courseRepo.On("DynamicFilterCourses", ctx, "rating", "4.0", "", 1, 10).Return(mockCourses, total, nil)

    response, err := suite.courseUseCase.FilterCourses(ctx, req)

    assert.NoError(suite.T(), err)
    assert.Len(suite.T(), response.Courses, 2)
    suite.courseRepo.AssertExpectations(suite.T())
}

func (suite *CourseUseCaseTestSuite) TestFilterCourses_Failure() {
    ctx := context.Background()
    category := "Programming Languages"
    req := FilterCoursesRequest{
        Category: &category,
        Page:     1,
        Limit:    10,
    }

    suite.courseRepo.On("DynamicFilterCourses", ctx, "category", "Programming Languages", "", 1, 10).Return([]schema.Course{}, 0, apierror.ErrInternalServer)

    response, err := suite.courseUseCase.FilterCourses(ctx, req)

    assert.Error(suite.T(), err)
    assert.Nil(suite.T(), response)
    suite.courseRepo.AssertExpectations(suite.T())
}

func TestCourseUseCaseTestSuite(t *testing.T) {
    suite.Run(t, new(CourseUseCaseTestSuite))
}

