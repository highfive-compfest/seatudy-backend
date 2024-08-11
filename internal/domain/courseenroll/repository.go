package courseenroll

import (
    "context"
    "github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
    "gorm.io/gorm"
)

type Repository interface {
    Create(ctx context.Context, enroll *schema.CourseEnroll) error
    GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]schema.CourseEnroll, error)
    GetByUserID(ctx context.Context, userID uuid.UUID) ([]schema.CourseEnroll, error)
	IsEnrolled(ctx context.Context, userID, courseID uuid.UUID) (bool, error)
}

type repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, enroll *schema.CourseEnroll) error {
    return r.db.WithContext(ctx).Create(enroll).Error
}

func (r *repository) GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]schema.CourseEnroll, error) {
    var enrolls []schema.CourseEnroll
    err := r.db.Where("course_id = ?", courseID).Find(&enrolls).Error
    return enrolls, err
}

func (r *repository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]schema.CourseEnroll, error) {
    var enrolls []schema.CourseEnroll
    err := r.db.Where("user_id = ?", userID).Find(&enrolls).Error
    return enrolls, err
}

func (r *repository) IsEnrolled(ctx context.Context, userID, courseID uuid.UUID) (bool, error) {
    var count int64
    err := r.db.Model(&schema.CourseEnroll{}).Where("user_id = ? AND course_id = ?", userID, courseID).Count(&count).Error
    if err != nil {
        return false, err
    }
    return count > 0, nil
}
