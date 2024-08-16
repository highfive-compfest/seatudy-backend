package course

import (
	"context"
	"log"
	"strconv"

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
	DynamicFilterCourses(ctx context.Context, filterType, filterValue, sort string, page, limit int) ([]schema.Course, int, error)
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

    // Use a raw SQL expression for the subquery to count enrollments
    enrollmentCountSQL := "(SELECT COUNT(*) FROM course_enrolls WHERE course_enrolls.course_id = courses.id)"

    // Query to fetch courses ordered by enrollment count and then by rating
    result := r.db.Model(&schema.Course{}).
        Select("courses.*, " + enrollmentCountSQL + " as enrollment_count").
        Order("enrollment_count DESC").
        Order("rating DESC").
        Preload("Materials.Attachments").
        Preload("Assignments.Attachments").
        Offset((page - 1) * pageSize).
        Limit(pageSize).
        Find(&courses)
    if result.Error != nil {
        return nil, 0, result.Error
    }

    var totalRecords int64
    r.db.Model(&schema.Course{}).Count(&totalRecords)
    return courses, int(totalRecords), nil
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

func (r *repository) DynamicFilterCourses(ctx context.Context, filterType, filterValue, sort string, page, limit int) ([]schema.Course, int, error) {
    var courses []schema.Course
    var total int64

    query := r.db.Model(&schema.Course{})

    switch filterType {
    case "category":
        query = query.Where("category = ?", filterValue)
    case "difficulty":
        query = query.Where("difficulty = ?", filterValue)
    case "rating":
		rating, _ := strconv.ParseFloat(filterValue, 32) // Assuming filterValue is a string; handle errors as needed
        upperBound := rating + 1.0
        query = query.Where("rating >= ? AND rating < ?", rating, upperBound)
    }

    if sort == "highest" {
        query = query.Order("rating DESC")
    } else if sort == "lowest" {
        query = query.Order("rating ASC")
    }

    // Count total results for pagination before limiting and offsetting
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Apply pagination
    query = query.Offset((page - 1) * limit).Limit(limit)

    // Execute the final query to retrieve the filtered courses
    if err := query.Find(&courses).Error; err != nil {
        return nil, 0, err
    }

    return courses, int(total) , nil
}
