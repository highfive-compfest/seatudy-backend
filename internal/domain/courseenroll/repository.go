package courseenroll

import (
    "context"
    "github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
    "gorm.io/gorm"
)

type Repository interface {
    Create(ctx context.Context, enroll *schema.CourseEnroll) error
    GetUsersByCourseID(ctx context.Context, courseID uuid.UUID) ([]schema.User, error)
    GetCoursesByUserID(ctx context.Context, userID uuid.UUID) ([]schema.Course, error)
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

func (r *repository) GetUsersByCourseID(ctx context.Context, courseID uuid.UUID) ([]schema.User, error) {
    var users []schema.User
    err := r.db.Joins("JOIN course_enrolls ON course_enrolls.user_id = users.id").
        Where("course_enrolls.course_id = ?", courseID).
        Find(&users).Error
    return users, err
}

func (r *repository) GetCoursesByUserID(ctx context.Context, userID uuid.UUID) ([]schema.Course, error) {
    var courses []schema.Course
    err := r.db.Joins("JOIN course_enrolls ON course_enrolls.course_id = courses.id").
        Where("course_enrolls.user_id = ?", userID).
        Find(&courses).Error
    return courses, err
}


func (r *repository) IsEnrolled(ctx context.Context, userID, courseID uuid.UUID) (bool, error) {
    var count int64
    err := r.db.Model(&schema.CourseEnroll{}).Where("user_id = ? AND course_id = ?", userID, courseID).Count(&count).Error
    if err != nil {
        return false, err
    }
    return count > 0, nil
}
