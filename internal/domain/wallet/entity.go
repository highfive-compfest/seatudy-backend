package wallet

import (
	"github.com/google/uuid"
	"time"
)

type Wallet struct {
	ID                   uuid.UUID             `gorm:"primaryKey"`
	UserID               uuid.UUID             `gorm:"not null;index"`
	Balance              int64                 `gorm:"not null; default:0; check:balance >= 0"`
	MidtransTransactions []MidtransTransaction `gorm:"foreignKey:WalletID"`
}

type MidtransStatus string

var (
	MidtransStatusChallenge MidtransStatus = "challenge"
	MidtransStatusSuccess   MidtransStatus = "success"
	MidtransStatusFailure   MidtransStatus = "failure"
	MidtransStatusPending   MidtransStatus = "pending"
)

type MidtransTransaction struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey"`
	WalletID  uuid.UUID      `json:"-" gorm:"not null;index"`
	Amount    int64          `json:"amount" gorm:"not null"`
	IsCredit  bool           `json:"-" gorm:"not null"`
	Status    MidtransStatus `json:"status" gorm:"type:midtrans_status;not null"`
	CreatedAt time.Time      `json:"created_at"`
	ExpireAt  time.Time      `json:"-"`
}
