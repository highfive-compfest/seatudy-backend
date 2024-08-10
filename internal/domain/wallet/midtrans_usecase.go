package wallet

import (
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type IMidtransUseCase interface {
	CreateTransaction(id string, amount int64) (*snap.Response, *midtrans.Error)
	VerifyPayment(notificationPayload map[string]any) error
}

type MidtransUseCase struct {
	walletUc *UseCase
}

func NewMidtransUseCase(walletUc *UseCase) *MidtransUseCase {
	return &MidtransUseCase{walletUc: walletUc}
}

func (muc *MidtransUseCase) CreateTransaction(id string, amount int64) (*snap.Response, *midtrans.Error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  id,
			GrossAmt: amount,
		},
		Expiry: &snap.ExpiryDetails{
			Unit:     "minute",
			Duration: 1,
		},
	}

	return snap.CreateTransaction(req)
}

func (muc *MidtransUseCase) VerifyPayment(notificationPayload map[string]any) error {
	// 3. Get order-id from payload
	orderId, exists := notificationPayload["order_id"].(string)
	if !exists {
		return apierror.ErrValidation
	}

	// 4. Check transaction to Midtrans with param orderId
	transactionStatusResp, e := coreapi.CheckTransaction(orderId)
	if e != nil {
		return nil // Return 200 for midtrans test notification, but do nothing
	}

	if transactionStatusResp != nil {
		transactionId, err := uuid.Parse(orderId)
		if err != nil {
			return apierror.ErrValidation
		}
		// 5. Do set transaction status based on response from check transaction status
		if transactionStatusResp.TransactionStatus == "capture" {
			if transactionStatusResp.FraudStatus == "challenge" {
				// set transaction status on your database to 'challenge'
				if err := muc.walletUc.VerifyPayment(transactionId, schema.MidtransStatusChallenge); err != nil {
					return apierror.ErrInternalServer
				}
				// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			} else if transactionStatusResp.FraudStatus == "accept" {
				// set transaction status on your database to 'success'
				if err := muc.walletUc.VerifyPayment(transactionId, schema.MidtransStatusSuccess); err != nil {
					return apierror.ErrInternalServer
				}
			}
		} else if transactionStatusResp.TransactionStatus == "settlement" {
			// set transaction status on your databaase to 'success'
			if err := muc.walletUc.VerifyPayment(transactionId, schema.MidtransStatusSuccess); err != nil {
				return apierror.ErrInternalServer
			}
		} else if transactionStatusResp.TransactionStatus == "deny" {
			// you can ignore 'deny', because most of the time it allows payment retries
			// and later can become success
		} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
			// set transaction status on your databaase to 'failure'
			if err := muc.walletUc.VerifyPayment(transactionId, schema.MidtransStatusFailure); err != nil {
				return apierror.ErrInternalServer
			}
		} else if transactionStatusResp.TransactionStatus == "pending" {
			// set transaction status on your databaase to 'pending' / waiting payment
			if err := muc.walletUc.VerifyPayment(transactionId, schema.MidtransStatusPending); err != nil {
				return apierror.ErrInternalServer
			}
		}
	}
	return nil
}
