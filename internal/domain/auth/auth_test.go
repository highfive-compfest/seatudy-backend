package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/config"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/user"
	"github.com/highfive-compfest/seatudy-backend/internal/jwtoken"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"os"
	"testing"
)

type MockAuthRepository struct {
	mock.Mock
}

type MockUserRepository struct {
	mock.Mock
}

type MockMailDialer struct {
	mock.Mock
}

func (m *MockAuthRepository) SaveOTP(ctx context.Context, email string, otp string) error {
	args := m.Called(ctx, email, otp)
	return args.Error(0)
}

func (m *MockAuthRepository) GetOTP(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}

func (m *MockAuthRepository) DeleteOTP(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockAuthRepository) SaveResetPasswordToken(ctx context.Context, email string, token string) error {
	args := m.Called(ctx, email, token)
	return args.Error(0)
}

func (m *MockAuthRepository) GetResetPasswordToken(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}

func (m *MockAuthRepository) DeleteResetPasswordToken(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockUserRepository) Create(user *schema.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*schema.User, error) {
	args := m.Called(id)
	userObj, ok := args.Get(0).(*schema.User)
	if !ok {
		return nil, args.Error(1)
	}
	return userObj, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*schema.User, error) {
	args := m.Called(email)
	userObj, ok := args.Get(0).(*schema.User)
	if !ok {
		return nil, args.Error(1)
	}
	return userObj, args.Error(1)
}

func (m *MockUserRepository) Update(user *schema.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateByEmail(email string, user *schema.User) error {
	args := m.Called(email, user)
	return args.Error(0)
}

func (m *MockMailDialer) DialAndSend(msg ...*gomail.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

type AuthUseCaseTestSuite struct {
	suite.Suite
	authRepo   *MockAuthRepository
	userRepo   *MockUserRepository
	mailDialer *MockMailDialer
	useCase    *UseCase
}

func (suite *AuthUseCaseTestSuite) SetupTest() {
	suite.authRepo = new(MockAuthRepository)
	suite.userRepo = new(MockUserRepository)
	suite.mailDialer = new(MockMailDialer)
	suite.useCase = NewUseCase(suite.authRepo, suite.userRepo, suite.mailDialer)
}

func (suite *AuthUseCaseTestSuite) TestRegister_Success() {
	req := &RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
		Role:     "student",
	}

	suite.userRepo.On("Create", mock.Anything).Return(nil)

	err := suite.useCase.Register(req)
	assert.NoError(suite.T(), err)
}

func (suite *AuthUseCaseTestSuite) TestRegister_EmailAlreadyRegistered() {
	req := &RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
		Role:     "student",
	}

	pgErr := pgconn.PgError{Code: "23505"}

	suite.userRepo.On("Create", mock.Anything).Return(&pgErr)

	err := suite.useCase.Register(req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrEmailAlreadyRegistered, err)
}

func (suite *AuthUseCaseTestSuite) TestRegister_InternalServerError() {
	req := &RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
		Role:     "student",
	}

	suite.userRepo.On("Create", mock.Anything).Return(gorm.ErrInvalidDB)

	err := suite.useCase.Register(req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInternalServer, err)
}

func (suite *AuthUseCaseTestSuite) TestLogin_Success() {
	_ = os.Setenv("ENV", "test")
	_ = os.Setenv("JWT_ACCESS_DURATION", "10m")
	config.LoadEnv()

	req := &LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	userObj := &schema.User{
		ID:              uuid.New(),
		Email:           req.Email,
		PasswordHash:    "$2a$10$JdUpBMZEt2gUid5JJU6XVuB6Mdiu.qs4r94cX28vB.Y7ovzg.PO9G",
		Name:            "Test User",
		Role:            "student",
		IsEmailVerified: true,
	}

	suite.userRepo.On("GetByEmail", req.Email).Return(userObj, nil)

	resp, err := suite.useCase.Login(req)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
}

func (suite *AuthUseCaseTestSuite) TestLogin_UserNotFound() {
	req := &LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	suite.userRepo.On("GetByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)

	resp, err := suite.useCase.Login(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), ErrInvalidCredentials, err)
}

func (suite *AuthUseCaseTestSuite) TestLogin_InvalidCredentials() {
	req := &LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	userObj := &schema.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: "$2a$10$7EqJtq98hPqEX7fNZaFWoOe5C5F1e4G1Z1Z1Z1Z1Z1Z1Z1Z1Z1Z1",
	}

	suite.userRepo.On("GetByEmail", req.Email).Return(userObj, nil)

	resp, err := suite.useCase.Login(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), ErrInvalidCredentials, err)
}

func (suite *AuthUseCaseTestSuite) TestLogin_InternalServerError() {
	req := &LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	suite.userRepo.On("GetByEmail", req.Email).Return(nil, gorm.ErrInvalidDB)

	resp, err := suite.useCase.Login(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), apierror.ErrInternalServer, err)
}

func (suite *AuthUseCaseTestSuite) TestRefresh_Success() {
	_ = os.Setenv("ENV", "test")
	_ = os.Setenv("JWT_REFRESH_DURATION", "720h")
	config.LoadEnv()

	token, _ := jwtoken.CreateRefreshJWT("01914b1c-4762-7d85-bc7a-6e81eda6f2c7")
	req := &RefreshRequest{
		RefreshToken: token,
	}

	suite.userRepo.On("GetByID", mock.Anything).Return(&schema.User{}, nil)

	resp, err := suite.useCase.Refresh(req)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
}

func (suite *AuthUseCaseTestSuite) TestRefresh_InvalidToken() {
	req := &RefreshRequest{
		RefreshToken: "invalid_refresh_token",
	}

	resp, err := suite.useCase.Refresh(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), apierror.ErrTokenInvalid, err)
}

func (suite *AuthUseCaseTestSuite) TestRefresh_TokenExpired() {
	_ = os.Setenv("ENV", "test")
	_ = os.Setenv("JWT_REFRESH_DURATION", "-720h")
	config.LoadEnv()

	token, _ := jwtoken.CreateRefreshJWT("01914b1c-4762-7d85-bc7a-6e81eda6f2c7")
	req := &RefreshRequest{
		RefreshToken: token,
	}

	resp, err := suite.useCase.Refresh(req)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), apierror.ErrTokenInvalid, err)
}

func (suite *AuthUseCaseTestSuite) TestSendOTP_Success() {
	_ = os.Setenv("ENV", "test")
	_ = os.Setenv("SMTP_EMAIL", "noreply@example.com")
	config.LoadEnv()

	ctx := context.WithValue(context.Background(), "user.email", "test@example.com")
	ctx = context.WithValue(ctx, "user.name", "Test User")
	ctx = context.WithValue(ctx, "user.is_email_verified", false)

	suite.authRepo.On("SaveOTP", ctx, "test@example.com", mock.Anything).Return(nil)
	suite.mailDialer.On("DialAndSend", mock.Anything).Return(nil)

	err := suite.useCase.SendOTP(ctx)
	assert.NoError(suite.T(), err)
}

func (suite *AuthUseCaseTestSuite) TestSendOTP_EmailAlreadyVerified() {
	_ = os.Setenv("ENV", "test")
	_ = os.Setenv("SMTP_EMAIL", "noreply@example.com")
	config.LoadEnv()

	ctx := context.WithValue(context.Background(), "user.email", "test@example.com")
	ctx = context.WithValue(ctx, "user.name", "Test User")
	ctx = context.WithValue(ctx, "user.is_email_verified", true)

	err := suite.useCase.SendOTP(ctx)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrEmailAlreadyVerified, err)
}

func (suite *AuthUseCaseTestSuite) TestVerifyOTP_Success() {
	ctx := context.WithValue(context.Background(), "user.email", "test@example.com")

	req := &VerifyEmailRequest{
		OTP: "123456",
	}

	suite.authRepo.On("GetOTP", ctx, "test@example.com").Return("123456", nil)
	suite.authRepo.On("DeleteOTP", ctx, "test@example.com").Return(nil)
	suite.userRepo.On("UpdateByEmail", "test@example.com", mock.Anything).Return(nil)

	err := suite.useCase.VerifyOTP(ctx, req)
	assert.NoError(suite.T(), err)
}

func (suite *AuthUseCaseTestSuite) TestVerifyOTP_ExpiredOTP() {
	ctx := context.WithValue(context.Background(), "user.email", "test@example.com")

	req := &VerifyEmailRequest{
		OTP: "123456",
	}

	suite.authRepo.On("GetOTP", ctx, "test@example.com").Return("", redis.Nil)

	err := suite.useCase.VerifyOTP(ctx, req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrExpiredOTP, err)
}

func (suite *AuthUseCaseTestSuite) TestVerifyOTP_InvalidOTP() {
	ctx := context.WithValue(context.Background(), "user.email", "test@example.com")

	req := &VerifyEmailRequest{
		OTP: "123456",
	}

	suite.authRepo.On("GetOTP", ctx, "test@example.com").Return("654321", nil)

	err := suite.useCase.VerifyOTP(ctx, req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidOTP, err)
}

func (suite *AuthUseCaseTestSuite) TestVerifyOTP_InternalServerError() {
	ctx := context.WithValue(context.Background(), "user.email", "test@example.com")

	req := &VerifyEmailRequest{
		OTP: "123456",
	}

	suite.authRepo.On("GetOTP", ctx, "test@example.com").Return("", redis.TxFailedErr)

	err := suite.useCase.VerifyOTP(ctx, req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInternalServer, err)
}

func (suite *AuthUseCaseTestSuite) TestSendResetPasswordLink_Success() {
	_ = os.Setenv("ENV", "test")
	_ = os.Setenv("FRONTEND_URL", "https://example.com")
	config.LoadEnv()

	req := &SendResetPasswordLinkRequest{
		Email: "test@example.com",
	}

	userObj := &schema.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	suite.userRepo.On("GetByEmail", req.Email).Return(userObj, nil)
	suite.authRepo.On("SaveResetPasswordToken", mock.Anything, req.Email, mock.Anything).Return(nil)
	suite.mailDialer.On("DialAndSend", mock.Anything).Return(nil)

	err := suite.useCase.SendResetPasswordLink(context.Background(), req)
	assert.NoError(suite.T(), err)
}

func (suite *AuthUseCaseTestSuite) TestSendResetPasswordLink_UserNotFound() {
	req := &SendResetPasswordLinkRequest{
		Email: "test@example.com",
	}

	suite.userRepo.On("GetByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)

	err := suite.useCase.SendResetPasswordLink(context.Background(), req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), user.ErrUserNotFound, err)
}

func (suite *AuthUseCaseTestSuite) TestSendResetPasswordLink_InternalServerError() {
	req := &SendResetPasswordLinkRequest{
		Email: "test@example.com",
	}

	suite.userRepo.On("GetByEmail", req.Email).Return(nil, gorm.ErrInvalidDB)

	err := suite.useCase.SendResetPasswordLink(context.Background(), req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInternalServer, err)
}

func (suite *AuthUseCaseTestSuite) TestResetPassword_Success() {
	req := &ResetPasswordRequest{
		Email:       "test@example.com",
		Token:       "valid_token",
		NewPassword: "newpassword123",
	}

	suite.authRepo.On("GetResetPasswordToken", mock.Anything, req.Email).Return("valid_token", nil)
	suite.userRepo.On("UpdateByEmail", req.Email, mock.Anything).Return(nil)
	suite.authRepo.On("DeleteResetPasswordToken", mock.Anything, req.Email).Return(nil)

	err := suite.useCase.ResetPassword(context.Background(), req)
	assert.NoError(suite.T(), err)
}

func (suite *AuthUseCaseTestSuite) TestResetPassword_ExpiredToken() {
	req := &ResetPasswordRequest{
		Email:       "test@example.com",
		Token:       "expired_token",
		NewPassword: "newpassword123",
	}

	suite.authRepo.On("GetResetPasswordToken", mock.Anything, req.Email).Return("", redis.Nil)

	err := suite.useCase.ResetPassword(context.Background(), req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrExpiredResetPasswordLink, err)
}

func (suite *AuthUseCaseTestSuite) TestResetPassword_InvalidToken() {
	req := &ResetPasswordRequest{
		Email:       "test@example.com",
		Token:       "invalid_token",
		NewPassword: "newpassword123",
	}

	suite.authRepo.On("GetResetPasswordToken", mock.Anything, req.Email).Return("valid_token", nil)

	err := suite.useCase.ResetPassword(context.Background(), req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidResetPasswordLink, err)
}

func (suite *AuthUseCaseTestSuite) TestResetPassword_InternalServerError() {
	req := &ResetPasswordRequest{
		Email: "test@example.com",
	}

	suite.authRepo.On("GetResetPasswordToken", mock.Anything, req.Email).Return("", redis.TxFailedErr)

	err := suite.useCase.ResetPassword(context.Background(), req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInternalServer, err)
}

func (suite *AuthUseCaseTestSuite) TestChangePassword_Success() {
	ctx := context.WithValue(context.Background(), "user.email", "test@example.com")

	req := &ChangePasswordRequest{
		OldPassword: "password123",
		NewPassword: "newpassword123",
	}

	userObj := &schema.User{
		Email:        "test@example.com",
		PasswordHash: "$2a$10$JdUpBMZEt2gUid5JJU6XVuB6Mdiu.qs4r94cX28vB.Y7ovzg.PO9G",
	}

	suite.userRepo.On("GetByEmail", "test@example.com").Return(userObj, nil)
	suite.userRepo.On("UpdateByEmail", "test@example.com", mock.Anything).Return(nil)

	err := suite.useCase.ChangePassword(ctx, req)
	assert.NoError(suite.T(), err)
}

func (suite *AuthUseCaseTestSuite) TestChangePassword_InvalidCredentials() {
	ctx := context.WithValue(context.Background(), "user.email", "test@example.com")

	req := &ChangePasswordRequest{
		OldPassword: "wrongpassword",
		NewPassword: "newpassword123",
	}

	userObj := &schema.User{
		Email:        "test@example.com",
		PasswordHash: "$2a$10$7EqJtq98hPqEX7fNZaFWoOe5C5F1e4G1Z1Z1Z1Z1Z1Z1Z1Z1Z1",
	}

	suite.userRepo.On("GetByEmail", "test@example.com").Return(userObj, nil)

	err := suite.useCase.ChangePassword(ctx, req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidCredentials, err)
}

func (suite *AuthUseCaseTestSuite) TestChangePassword_InternalServerError() {
	ctx := context.WithValue(context.Background(), "user.email", "test@example.com")

	req := &ChangePasswordRequest{
		OldPassword: "password123",
		NewPassword: "newpassword123",
	}

	suite.userRepo.On("GetByEmail", "test@example.com").Return(nil, gorm.ErrInvalidDB)

	err := suite.useCase.ChangePassword(ctx, req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInternalServer, err)
}

func TestAuthUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(AuthUseCaseTestSuite))
}
