package submission

import (
	"context"

	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, s *schema.Submission) error
	Update(ctx context.Context, s *schema.Submission) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*schema.Submission, error)
	GetAllByAssignment(ctx context.Context, assignmentID uuid.UUID) ([]schema.Submission, error)
	CheckSubmissionExists(ctx context.Context, userID, assignmentID uuid.UUID) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, s *schema.Submission) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *repository) Update(ctx context.Context, s *schema.Submission) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&schema.Submission{}, id).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*schema.Submission, error) {
	var submission schema.Submission
	if err := r.db.Preload("Attachments").First(&submission, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}

func (r *repository) GetAllByAssignment(ctx context.Context, assignmentID uuid.UUID) ([]schema.Submission, error) {
	var submissions []schema.Submission
	if err := r.db.Where("assignment_id = ?", assignmentID).Preload("Attachments").Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *repository) CheckSubmissionExists(ctx context.Context, userID, assignmentID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&schema.Submission{}).
		Where("user_id = ? AND assignment_id = ? AND deleted_at IS NULL", userID, assignmentID).
		Count(&count).Error
	return count > 0, err
}