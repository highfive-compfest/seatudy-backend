// file: attachment/repository.go

package attachment

import (
	"context"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, att *schema.Attachment) error
	Update(ctx context.Context, att *schema.Attachment) error
	GetByID(ctx context.Context, id uuid.UUID) (*schema.Attachment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, att *schema.Attachment) error {
	return r.db.WithContext(ctx).Create(att).Error
}

func (r *repository) Update(ctx context.Context, att *schema.Attachment) error {
	return r.db.WithContext(ctx).Save(att).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*schema.Attachment, error) {
	var att schema.Attachment
	result := r.db.WithContext(ctx).First(&att, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &att, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.Delete(&schema.Attachment{}, id).Error
}
