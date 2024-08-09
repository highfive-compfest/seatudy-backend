package material


import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, mat *Material) error
	GetByID(ctx context.Context, id uuid.UUID) (*Material, error)
	GetAll(ctx context.Context) ([]*Material, error)
	Update(ctx context.Context, mat *Material) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, mat *Material) error {
	return r.db.WithContext(ctx).Create(mat).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Material, error) {
	var mat Material
	if err := r.db.Preload("Attachments").First(&mat, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &mat, nil
}

func (r *repository) GetAll(ctx context.Context) ([]*Material, error) {
    var materials []*Material
    result := r.db.Preload("Attachments").Find(&materials)  // Preload the attachments
    if result.Error != nil {
        return nil, result.Error
    }
    return materials, nil
}

func (r *repository) Update(ctx context.Context, mat *Material) error {
	return r.db.WithContext(ctx).Save(mat).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Material{}, id).Error
}