package wallet

import (
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"gorm.io/gorm"
)

type IRepository interface {
	Create(tx *gorm.DB, wallet *schema.Wallet) error
	CreateMidtransTransaction(tx *gorm.DB, transaction *schema.MidtransTransaction) error

	GetByUserID(tx *gorm.DB, userID uuid.UUID) (*schema.Wallet, error)
	GetMidtransTransactionByID(tx *gorm.DB, transactionID uuid.UUID) (*schema.MidtransTransaction, error)
	GetMidtransTransactionsByWalletID(tx *gorm.DB, walletID uuid.UUID, isCredit bool, page,
		limit int) ([]*schema.MidtransTransaction, int64, error)

	UpdateMidtransTransaction(tx *gorm.DB, transaction *schema.MidtransTransaction) error

	TopUpSuccess(transactionID uuid.UUID) error
	TransferByUserID(tx *gorm.DB, fromUserID, toUserID uuid.UUID, amount int64) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{db}
}

func (r *Repository) Create(tx *gorm.DB, wallet *schema.Wallet) error {
	if tx == nil {
		tx = r.db
	}

	return tx.Create(wallet).Error
}

func (r *Repository) CreateMidtransTransaction(tx *gorm.DB, transaction *schema.MidtransTransaction) error {
	if tx == nil {
		tx = r.db
	}

	return tx.Create(transaction).Error
}

func (r *Repository) GetByUserID(tx *gorm.DB, userID uuid.UUID) (*schema.Wallet, error) {
	if tx == nil {
		tx = r.db
	}

	var wallet schema.Wallet
	tx = tx.Where("user_id = ?", userID).First(&wallet)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &wallet, nil
}

func (r *Repository) GetMidtransTransactionByID(tx *gorm.DB, transactionID uuid.UUID) (*schema.MidtransTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	var transaction schema.MidtransTransaction
	tx = tx.Where("id = ?", transactionID).First(&transaction)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &transaction, nil
}

func (r *Repository) GetMidtransTransactionsByWalletID(tx *gorm.DB, walletID uuid.UUID, isCredit bool, page, limit int) ([]*schema.MidtransTransaction, int64, error) {
	if tx == nil {
		tx = r.db
	}

	var transactions []*schema.MidtransTransaction
	tx = tx.Model(&schema.MidtransTransaction{}).Where("wallet_id = ? AND is_credit = ?", walletID, isCredit)
	var total int64
	tx.Count(&total)
	tx.Order("id DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&transactions)
	if tx.Error != nil {
		return nil, 0, tx.Error
	}
	return transactions, total, nil
}

func (r *Repository) UpdateMidtransTransaction(tx *gorm.DB, transaction *schema.MidtransTransaction) error {
	if tx == nil {
		tx = r.db
	}
	tx.Updates(transaction)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *Repository) TopUpSuccess(transactionID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		transaction, err := r.GetMidtransTransactionByID(tx, transactionID)
		if err != nil {
			return err
		}

		if err := r.UpdateMidtransTransaction(tx,
			&schema.MidtransTransaction{ID: transactionID, Status: schema.MidtransStatusSuccess}); err != nil {
			return err
		}

		return tx.Model(&schema.Wallet{}).Where("id = ?", transaction.WalletID).
			Update("balance", gorm.Expr("balance + ?", transaction.Amount)).Error
	})
}

func (r *Repository) TransferByUserID(tx *gorm.DB, fromUserID, toUserID uuid.UUID, amount int64) error {
	if tx == nil {
		tx = r.db
	}

	return tx.Transaction(func(tx *gorm.DB) error {
		fromWallet, err := r.GetByUserID(tx, fromUserID)
		if err != nil {
			return err
		}

		toWallet, err := r.GetByUserID(tx, toUserID)
		if err != nil {
			return err
		}

		// debit sender
		if fromWallet.Balance < amount {
			return apierror.ErrInsufficientBalance
		}
		if err := tx.Model(fromWallet).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		// credit receiver
		if err := tx.Model(toWallet).Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		return nil
	})
}
