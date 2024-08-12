package assignment

import (
	"context"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, a *schema.Assignment) error
	Update(ctx context.Context, a *schema.Assignment) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*schema.Assignment, error)
	GetByCourseID(ctx context.Context, courseId uuid.UUID) ([]*schema.Assignment, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, a *schema.Assignment) error {
	return r.db.Create(a).Error
}

func (r *repository) Update(ctx context.Context, a *schema.Assignment) error {
	return r.db.Save(a).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.Delete(&schema.Assignment{}, id).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*schema.Assignment, error) {
	var assignment schema.Assignment
	result := r.db.Preload("Attachments").Where("id = ?", id).First(&assignment)
	return &assignment, result.Error
}

func (r *repository) GetByCourseID(ctx context.Context, courseId uuid.UUID) ([]*schema.Assignment, error) {
	var assignments []*schema.Assignment
	result := r.db.Preload("Attachments").Where("course_id = ?", courseId).Find(&assignments)
	return assignments, result.Error
}
