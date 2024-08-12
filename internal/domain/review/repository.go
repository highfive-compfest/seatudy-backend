package review

import (
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"gorm.io/gorm"
)

type IRepository interface {
	Create(review *schema.Review, newCourseRating float32, newCourseRatingCount int64) error
	GetByID(id uuid.UUID) (*schema.Review, error)
	Get(condition map[string]any, page int, limit int) ([]schema.Review, int64, error)
	Update(review *schema.Review, courseID uuid.UUID, newCourseRating float32) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	return &repository{db}
}

func (r *repository) Create(review *schema.Review, newCourseRating float32, newCourseReviewCount int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&schema.Course{ID: review.CourseID}).Updates(map[string]any{
			"rating":       newCourseRating,
			"review_count": newCourseReviewCount,
		}).Error; err != nil {
			return err
		}

		return r.db.Create(review).Error
	})
}

func (r *repository) GetByID(id uuid.UUID) (*schema.Review, error) {
	var review schema.Review
	tx := r.db.First(&review, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &review, nil
}

func (r *repository) Get(condition map[string]any, page int, limit int) ([]schema.Review, int64, error) {
	var reviews []schema.Review
	var total int64

	tx := r.db.Model(&schema.Review{})

	for key, value := range condition {
		tx = tx.Where(key+" = ?", value)
	}

	tx.Count(&total)

	tx.Order("id DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&reviews)

	return reviews, total, tx.Error
}

func (r *repository) Update(review *schema.Review, courseID uuid.UUID, newCourseRating float32) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if newCourseRating != 0 {
			if err := tx.Model(&schema.Course{ID: courseID}).UpdateColumn("rating", newCourseRating).Error; err != nil {
				return err
			}
		}
		tx.Updates(review)
		if tx.Error != nil {
			return tx.Error
		}
		return nil
	})
}

func (r *repository) Delete(id uuid.UUID) error {
	tx := r.db.Delete(&schema.Review{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
