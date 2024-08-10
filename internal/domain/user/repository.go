package user

import (
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/wallet"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"gorm.io/gorm"
)

type IRepository interface {
	Create(user *schema.User) error
	GetByID(id uuid.UUID) (*schema.User, error)
	GetByEmail(email string) (*schema.User, error)
	Update(user *schema.User) error
	UpdateByEmail(email string, user *schema.User) error
}

type repository struct {
	db         *gorm.DB
	walletRepo wallet.IRepository
}

func NewRepository(db *gorm.DB, walletRepo wallet.IRepository) IRepository {
	return &repository{db, walletRepo}
}

func (r *repository) Create(user *schema.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		walletID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		return r.walletRepo.Create(tx, &schema.Wallet{
			ID:     walletID,
			UserID: user.ID,
		})
	})
}

func (r *repository) GetByID(id uuid.UUID) (*schema.User, error) {
	var user schema.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByEmail(email string) (*schema.User, error) {
	var user schema.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) Update(user *schema.User) error {
	tx := r.db.Updates(user)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) UpdateByEmail(email string, user *schema.User) error {
	tx := r.db.Where("email = ?", email).Updates(user)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
