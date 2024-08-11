package course

import (
	"context"

	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context, page, pageSize int) ([]schema.Course, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (schema.Course, error)
	GetRating(ctx context.Context, courseID uuid.UUID) (float32, int64, error)
	Create(ctx context.Context, course *schema.Course) error
	Update(ctx context.Context, course *schema.Course) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByInstructorID(ctx context.Context, instructorID uuid.UUID, page, pageSize int) ([]schema.Course, int, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context, page, pageSize int) ([]schema.Course, int, error) {
    var courses []schema.Course
    result := r.db.Preload("Materials.Attachments").Offset((page - 1) * pageSize).Limit(pageSize).Find(&courses)
    if result.Error != nil {
        return nil, 0, result.Error
    }
    var totalRecords int64
    r.db.Model(&schema.Course{}).Count(&totalRecords)
    return courses, int(totalRecords), nil
}

func (r *repository) FindByInstructorID(ctx context.Context, instructorID uuid.UUID, page, pageSize int) ([]schema.Course, int, error) {
    var courses []schema.Course
    result := r.db.Where("instructor_id = ?", instructorID).Offset((page - 1) * pageSize).Limit(pageSize).Find(&courses)
    if result.Error != nil {
        return nil, 0, result.Error
    }
    var totalRecords int64
    r.db.Model(&schema.Course{}).Where("instructor_id = ?", instructorID).Count(&totalRecords)
    return courses, int(totalRecords), nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (schema.Course, error) {
	var course schema.Course
	if err := r.db.Preload("Materials.Attachments").First(&course, "id = ?", id).Error; err != nil {
		return schema.Course{}, err
	}
	return course, nil
}

func (r *repository) GetRating(ctx context.Context, courseID uuid.UUID) (float32, int64, error) {
	var course schema.Course
	if err := r.db.Select("rating", "review_count").First(&course, "id = ?", courseID).Error; err != nil {
		return 0, 0, err
	}
	return course.Rating, course.ReviewCount, nil
}

func (r *repository) Create(ctx context.Context, course *schema.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

func (r *repository) Update(ctx context.Context, course *schema.Course) error {
	return r.db.WithContext(ctx).Save(course).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&schema.Course{}, id).Error
}

