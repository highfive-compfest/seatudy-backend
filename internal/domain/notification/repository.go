package notification

import (
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"gorm.io/gorm"
)

type IRepository interface {
	Create(notification *schema.Notification) error
	GetByUserID(userID uuid.UUID, limit, offset int) ([]*schema.Notification, int64, error)
	GetUnreadCount(userID uuid.UUID) (int64, error)
	UpdateRead(notificationID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	return &repository{db: db}
}

func (r *repository) Create(notification *schema.Notification) error {
	return r.db.Create(notification).Error
}

func (r *repository) GetByUserID(userID uuid.UUID, limit, offset int) ([]*schema.Notification, int64, error) {
	var notifications []*schema.Notification
	var total int64

	tx := r.db.Model(&schema.Notification{}).Where("user_id = ?", userID)

	tx.Count(&total)

	tx.Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications)

	return notifications, total, tx.Error
}

func (r *repository) GetUnreadCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&schema.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

func (r *repository) UpdateRead(notificationID uuid.UUID) error {
	return r.db.Model(&schema.Notification{}).Where("id = ?", notificationID).Update("is_read", true).Error
}
