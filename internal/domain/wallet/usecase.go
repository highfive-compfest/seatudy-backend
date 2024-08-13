package wallet

import (
	"context"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/pagination"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"log"
	"time"
)

type UseCase struct {
	repo   IRepository
	MidtUc IMidtransUseCase
}

func NewUseCase(repo IRepository, midtUc IMidtransUseCase) *UseCase {
	return &UseCase{repo: repo, MidtUc: midtUc}
}

func (uc *UseCase) TopUp(ctx context.Context, req *TopUpRequest) (*TopUpResponse, error) {
	// Get user id from context
	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return nil, err
	}

	// 2. Get user wallet by user id
	wallet, err := uc.repo.GetByUserID(nil, userID)
	if err != nil {
		return nil, err
	}

	// 3. Create midtrans transaction
	transactionID, err := uuid.NewV7()
	if err != nil {
		return nil, apierror.ErrInternalServer
	}
	transaction := &schema.MidtransTransaction{
		ID:       transactionID,
		WalletID: wallet.ID,
		Amount:   req.Amount,
		IsCredit: true,
		Status:   schema.MidtransStatusPending,
	}

	// 4. Create midtrans transaction in database
	if err := uc.repo.CreateMidtransTransaction(nil, transaction); err != nil {
		return nil, err
	}

	// 5. Create transaction in midtrans
	snapResp, midtErr := uc.MidtUc.CreateTransaction(transaction.ID.String(), req.Amount)
	if midtErr != nil {
		return nil, midtErr
	}

	return &TopUpResponse{RedirectURL: snapResp.RedirectURL}, nil
}

func (uc *UseCase) VerifyPayment(transactionID uuid.UUID, status schema.MidtransStatus) error {
	if status == schema.MidtransStatusSuccess {
		if err := uc.repo.TopUpSuccess(transactionID); err != nil {
			log.Println("Error top up success: ", err)
			return apierror.ErrInternalServer
		}
	} else {
		if err := uc.repo.UpdateMidtransTransaction(nil,
			&schema.MidtransTransaction{ID: transactionID, Status: status}); err != nil {
			log.Println("Error update midtrans transaction status: ", err)
			return apierror.ErrInternalServer
		}
	}

	return nil
}

func (uc *UseCase) GetBalance(ctx context.Context) (*GetBalanceResponse, error) {
	// Get user id from context
	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return nil, apierror.ErrTokenInvalid
	}

	// Get user wallet by user id
	wallet, err := uc.repo.GetByUserID(nil, userID)
	if err != nil {
		log.Println("Error get wallet by user id: ", err)
		return nil, apierror.ErrInternalServer
	}

	return &GetBalanceResponse{Balance: wallet.Balance}, nil
}

func (uc *UseCase) GetMidtransTransactionsByUser(ctx context.Context,
	req *GetMidtransTransactionsRequest) (*pagination.GetResourcePaginatedResponse, error) {
	// Get user id from context
	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return nil, apierror.ErrTokenInvalid
	}

	isCredit := ctx.Value("user.role").(string) == "student"

	// Get user wallet by user id
	wallet, err := uc.repo.GetByUserID(nil, userID)
	if err != nil {
		log.Println("Error get wallet by user id: ", err)
		return nil, apierror.ErrInternalServer
	}

	// Get midtrans transactions by wallet id
	midtransTransactions, total, err := uc.repo.GetMidtransTransactionsByWalletID(nil, wallet.ID, isCredit, req.Page, req.Limit)
	if err != nil {
		log.Println("Error get midtrans transactions by wallet id: ", err)
		return nil, apierror.ErrInternalServer
	}

	for _, v := range midtransTransactions {
		if v.Status == schema.MidtransStatusPending && v.ExpireAt.Before(time.Now()) {
			v.Status = schema.MidtransStatusFailure
		}
	}

	resp := pagination.GetResourcePaginatedResponse{
		Data:       midtransTransactions,
		Pagination: pagination.NewPagination(int(total), req.Page, req.Limit),
	}

	return &resp, nil
}
