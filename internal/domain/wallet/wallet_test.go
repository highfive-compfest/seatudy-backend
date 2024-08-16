package wallet

import (
	"context"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type MockMidtransUseCase struct {
	mock.Mock
}

func (m *MockMidtransUseCase) CreateTransaction(id string, amount int64) (*snap.Response, *midtrans.Error) {
	args := m.Called(id, amount)
	var snapResp *snap.Response
	if args.Get(0) != nil {
		snapResp = args.Get(0).(*snap.Response)
	}
	var midtransErr *midtrans.Error
	if args.Get(1) != nil {
		midtransErr = args.Get(1).(*midtrans.Error)
	}
	return snapResp, midtransErr
}

func (m *MockMidtransUseCase) VerifyPayment(notificationPayload map[string]any) error {
	args := m.Called(notificationPayload)
	return args.Error(0)
}

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(tx *gorm.DB, wallet *schema.Wallet) error {
	args := m.Called(tx, wallet)
	return args.Error(0)
}

func (m *MockRepository) CreateMidtransTransaction(tx *gorm.DB, transaction *schema.MidtransTransaction) error {
	args := m.Called(tx, transaction)
	return args.Error(0)
}

func (m *MockRepository) GetByUserID(tx *gorm.DB, userID uuid.UUID) (*schema.Wallet, error) {
	args := m.Called(tx, userID)
	var wallet *schema.Wallet
	if args.Get(0) != nil {
		wallet = args.Get(0).(*schema.Wallet)
	}
	return wallet, args.Error(1)
}

func (m *MockRepository) GetMidtransTransactionByID(tx *gorm.DB, transactionID uuid.UUID) (*schema.MidtransTransaction, error) {
	args := m.Called(tx, transactionID)
	return args.Get(0).(*schema.MidtransTransaction), args.Error(1)
}

func (m *MockRepository) GetMidtransTransactionsByWalletID(tx *gorm.DB, walletID uuid.UUID, isCredit bool, page, limit int) ([]*schema.MidtransTransaction, int64, error) {
	args := m.Called(tx, walletID, isCredit, page, limit)
	var transactions []*schema.MidtransTransaction

	if args.Get(0) != nil {
		transactions = args.Get(0).([]*schema.MidtransTransaction)
	}

	return transactions, args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) UpdateMidtransTransaction(tx *gorm.DB, transaction *schema.MidtransTransaction) error {
	args := m.Called(tx, transaction)
	return args.Error(0)
}

func (m *MockRepository) TopUpSuccess(transactionID uuid.UUID) error {
	args := m.Called(transactionID)
	return args.Error(0)
}

func (m *MockRepository) TransferByUserID(tx *gorm.DB, fromUserID, toUserID uuid.UUID, amount int64) error {
	args := m.Called(tx, fromUserID, toUserID, amount)
	return args.Error(0)
}

type WalletUseCaseTestSuite struct {
	suite.Suite
	repo   *MockRepository
	midtUc *MockMidtransUseCase
	uc     *UseCase
}

func (suite *WalletUseCaseTestSuite) SetupTest() {
	suite.repo = new(MockRepository)
	suite.midtUc = new(MockMidtransUseCase)
	suite.uc = NewUseCase(suite.repo, suite.midtUc)
}

func (suite *WalletUseCaseTestSuite) TestTopUp_Success() {
	ctx := context.WithValue(context.Background(), "user.id", "123e4567-e89b-12d3-a456-426614174000")
	req := &TopUpRequest{Amount: 10000}
	wallet := &schema.Wallet{ID: uuid.New(), UserID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), Balance: 0}

	suite.repo.On("GetByUserID", mock.Anything, wallet.UserID).Return(wallet, nil)
	suite.repo.On("CreateMidtransTransaction", mock.Anything, mock.Anything).Return(nil)
	suite.midtUc.On("CreateTransaction", mock.AnythingOfType("string"), req.Amount).Return(&snap.Response{RedirectURL: "http://example.com"}, nil)

	res, err := suite.uc.TopUp(ctx, req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "http://example.com", res.RedirectURL)
}

func (suite *WalletUseCaseTestSuite) TestTopUp_InvalidUserID() {
	ctx := context.WithValue(context.Background(), "user.id", "invalid-uuid")
	req := &TopUpRequest{Amount: 10000}

	_, err := suite.uc.TopUp(ctx, req)
	assert.Error(suite.T(), err)
}

func (suite *WalletUseCaseTestSuite) TestTopUp_RepoError() {
	ctx := context.WithValue(context.Background(), "user.id", "123e4567-e89b-12d3-a456-426614174000")
	req := &TopUpRequest{Amount: 10000}
	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	suite.repo.On("GetByUserID", mock.Anything, userID).Return(nil, gorm.ErrInvalidDB)

	_, err := suite.uc.TopUp(ctx, req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInternalServer.Build(), err)
}

func (suite *WalletUseCaseTestSuite) TestTopUp_MidtransError() {
	ctx := context.WithValue(context.Background(), "user.id", "123e4567-e89b-12d3-a456-426614174000")
	req := &TopUpRequest{Amount: 10000}
	wallet := &schema.Wallet{ID: uuid.New(), UserID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), Balance: 0}

	suite.repo.On("GetByUserID", mock.Anything, wallet.UserID).Return(wallet, nil)
	suite.repo.On("CreateMidtransTransaction", mock.Anything, mock.Anything).Return(nil)
	suite.midtUc.On("CreateTransaction", mock.AnythingOfType("string"), req.Amount).Return(nil, &midtrans.Error{})

	_, err := suite.uc.TopUp(ctx, req)
	assert.Error(suite.T(), err)
}

func (suite *WalletUseCaseTestSuite) TestVerifyPayment_Success() {
	transactionID := uuid.New()

	suite.repo.On("TopUpSuccess", transactionID).Return(nil)
	err := suite.uc.VerifyPayment(transactionID, schema.MidtransStatusSuccess)
	assert.NoError(suite.T(), err)
}

func (suite *WalletUseCaseTestSuite) TestVerifyPayment_RepoError() {
	transactionID := uuid.New()

	suite.repo.On("TopUpSuccess", transactionID).Return(gorm.ErrInvalidDB)
	err := suite.uc.VerifyPayment(transactionID, schema.MidtransStatusSuccess)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInternalServer.Build(), err)
}

func (suite *WalletUseCaseTestSuite) TestGetBalance_Success() {
	ctx := context.WithValue(context.Background(), "user.id", "123e4567-e89b-12d3-a456-426614174000")
	wallet := &schema.Wallet{ID: uuid.New(), UserID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), Balance: 10000}

	suite.repo.On("GetByUserID", mock.Anything, wallet.UserID).Return(wallet, nil)

	res, err := suite.uc.GetBalance(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), wallet.Balance, res.Balance)
}

func (suite *WalletUseCaseTestSuite) TestGetBalance_InvalidUserID() {
	ctx := context.WithValue(context.Background(), "user.id", "invalid-uuid")

	_, err := suite.uc.GetBalance(ctx)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrTokenInvalid.Build(), err)
}

func (suite *WalletUseCaseTestSuite) TestGetBalance_RepoError() {
	ctx := context.WithValue(context.Background(), "user.id", "123e4567-e89b-12d3-a456-426614174000")
	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	suite.repo.On("GetByUserID", mock.Anything, userID).Return(nil, gorm.ErrInvalidDB)

	_, err := suite.uc.GetBalance(ctx)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInternalServer.Build(), err)
}

func (suite *WalletUseCaseTestSuite) TestGetMidtransTransactionsByUser_Success() {
	ctx := context.WithValue(context.Background(), "user.id", "123e4567-e89b-12d3-a456-426614174000")
	ctx = context.WithValue(ctx, "user.role", "student")
	req := &GetMidtransTransactionsRequest{Page: 1, Limit: 10}
	wallet := &schema.Wallet{ID: uuid.New(), UserID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), Balance: 10000}
	transactions := []*schema.MidtransTransaction{
		{ID: uuid.New(), WalletID: wallet.ID, Amount: 10000, Status: schema.MidtransStatusSuccess, CreatedAt: time.Now()},
	}

	suite.repo.On("GetByUserID", mock.Anything, wallet.UserID).Return(wallet, nil)
	suite.repo.On("GetMidtransTransactionsByWalletID", mock.Anything, wallet.ID, true, req.Page, req.Limit).Return(transactions, int64(len(transactions)), nil)

	res, err := suite.uc.GetMidtransTransactionsByUser(ctx, req)
	assert.NoError(suite.T(), err)
	resData, _ := res.Data.([]*schema.MidtransTransaction)
	assert.Equal(suite.T(), len(transactions), len(resData))
}

func (suite *WalletUseCaseTestSuite) TestGetMidtransTransactionsByUser_InvalidUserID() {
	ctx := context.WithValue(context.Background(), "user.id", "invalid-uuid")
	req := &GetMidtransTransactionsRequest{Page: 1, Limit: 10}

	_, err := suite.uc.GetMidtransTransactionsByUser(ctx, req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrTokenInvalid.Build(), err)
}

func (suite *WalletUseCaseTestSuite) TestGetMidtransTransactionsByUser_RepoError() {
	ctx := context.WithValue(context.Background(), "user.id", "123e4567-e89b-12d3-a456-426614174000")
	ctx = context.WithValue(ctx, "user.role", "student")
	req := &GetMidtransTransactionsRequest{Page: 1, Limit: 10}
	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	wallet := &schema.Wallet{ID: uuid.New(), UserID: userID, Balance: 10000}

	suite.repo.On("GetByUserID", mock.Anything, userID).Return(wallet, nil)
	suite.repo.On("GetMidtransTransactionsByWalletID", mock.Anything, wallet.ID, true, req.Page, req.Limit).
		Return(nil, int64(0), apierror.ErrInternalServer.Build())

	_, err := suite.uc.GetMidtransTransactionsByUser(ctx, req)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), apierror.ErrInternalServer.Build(), err)
}

func TestWalletUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(WalletUseCaseTestSuite))
}
