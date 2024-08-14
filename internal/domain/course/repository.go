package course

import (
	"context"
	"log"

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
	FindByPopularity(ctx context.Context, page, pageSize int) ([]schema.Course, int, error)
	GetUserCourseProgress(ctx context.Context, courseID, userID uuid.UUID) (float64, error)
	SearchByTitle(ctx context.Context, title string, page, pageSize int) ([]schema.Course, int, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context, page, pageSize int) ([]schema.Course, int, error) {
    var courses []schema.Course
    result := r.db.Preload("Materials.Attachments").Preload("Assignments.Attachments").Offset((page - 1) * pageSize).Limit(pageSize).Find(&courses)
    if result.Error != nil {
        return nil, 0, result.Error
    }
    var totalRecords int64
    r.db.Model(&schema.Course{}).Count(&totalRecords)
    return courses, int(totalRecords), nil
}

func (r *repository) FindByPopularity(ctx context.Context, page, pageSize int) ([]schema.Course, int, error) {
    var courses []schema.Course
    result := r.db.Order("rating DESC").Preload("Materials.Attachments").Preload("Assignments.Attachments").Offset((page - 1) * pageSize).Limit(pageSize).Find(&courses)
    if result.Error != nil {
        return nil,0, result.Error
    }

	var totalRecords int64
    r.db.Model(&schema.Course{}).Count(&totalRecords)
    return courses,int(totalRecords), nil
}

func (r *repository) GetUserCourseProgress(ctx context.Context, courseID, userID uuid.UUID) (float64, error) {
	var totalAssignments, completedAssignments int64


	if err := r.db.Model(&schema.Assignment{}).Where("course_id = ?", courseID).Count(&totalAssignments).Error; err != nil {
		return 0, err
	}


	if err := r.db.Model(&schema.Assignment{}).
		Joins("inner join submissions on submissions.assignment_id = assignments.id").
		Where("assignments.course_id = ? AND submissions.user_id = ? AND submissions.deleted_at IS NULL", courseID, userID).
		Count(&completedAssignments).Error; err != nil {
		return 0, err
		}
	
		log.Println(totalAssignments)
		log.Println(completedAssignments)
	// Calculate progress as a percentage
	var progress float64
	if totalAssignments > 0 {
		progress = (float64(completedAssignments) / float64(totalAssignments)) * 100
	}

	return progress, nil
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

func (r *repository) SearchByTitle(ctx context.Context, title string, page, pageSize int) ([]schema.Course, int, error) {
    var courses []schema.Course
    result := r.db.Where("title ILIKE ?", "%"+title+"%").Offset((page - 1) * pageSize).Limit(pageSize).Find(&courses)
    if result.Error != nil {
        return nil, 0, result.Error
    }
    var totalRecords int64
    r.db.Model(&schema.Course{}).Where("title ILIKE ?", "%"+title+"%").Count(&totalRecords)
    return courses, int(totalRecords), nil
}

