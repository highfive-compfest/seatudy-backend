package review

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/course"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/courseenroll"
	"github.com/highfive-compfest/seatudy-backend/internal/pagination"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"log"
)

type UseCase struct {
	repo       IRepository
	courseRepo course.Repository
	enrollUc   *courseenroll.UseCase
}

func NewUseCase(repo IRepository, courseRepo course.Repository, enrollUc *courseenroll.UseCase) *UseCase {
	return &UseCase{repo: repo, courseRepo: courseRepo, enrollUc: enrollUc}
}

func (uc *UseCase) Create(ctx context.Context, req *CreateReviewRequest) (*CreateReviewResponse, error) {
	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return nil, apierror.ErrTokenInvalid
	}

	// Check if user is enrolled in the course
	ok, err := uc.enrollUc.CheckEnrollment(ctx, userID, req.CourseID)
	if err != nil {
		log.Println("Error checking enrollment: ", err)
		return nil, apierror.ErrInternalServer
	}
	if !ok {
		return nil, courseenroll.ErrNotEnrolled
	}

	reviewID, err := uuid.NewV7()
	if err != nil {
		return nil, apierror.ErrInternalServer
	}

	rating, count, err := uc.courseRepo.GetRating(ctx, req.CourseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, course.ErrCourseNotFound
		}
		log.Print("Error getting course rating: ", err)
		return nil, apierror.ErrInternalServer
	}

	// mean_{n+1} = \frac{mean_n * n + x_{new}}{n+1}
	newRating := (rating*float32(count) + float32(req.Rating)) / float32(count+1)
	newCount := count + 1

	review := &schema.Review{
		ID:       reviewID,
		UserID:   userID,
		CourseID: req.CourseID,
		Rating:   req.Rating,
		Feedback: req.Feedback,
	}

	err = uc.repo.Create(review, newRating, newCount)
	if err != nil {
		var pgErr *pgconn.PgError
		ok := errors.As(err, &pgErr)
		if ok && pgErr.Code == "23505" {
			return nil, ErrCourseAlreadyReviewed
		}
		log.Print("Error creating review: ", err)
		return nil, apierror.ErrInternalServer
	}

	return &CreateReviewResponse{
		ID: reviewID,
	}, nil
}

func (uc *UseCase) Get(ctx context.Context, req *GetReviewsRequest) (*pagination.GetResourcePaginatedResponse, error) {
	conditions := make(map[string]any)
	if req.CourseID != "" {
		conditions["course_id"] = req.CourseID
	}
	if req.Rating != 0 {
		conditions["rating"] = req.Rating
	}

	reviews, total, err := uc.repo.Get(conditions, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	return &pagination.GetResourcePaginatedResponse{
		Data:       reviews,
		Pagination: pagination.NewPagination(int(total), req.Page, req.Limit),
	}, nil

}

func (uc *UseCase) Update(ctx context.Context, req *UpdateReviewRequest) error {
	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return apierror.ErrInternalServer
	}

	reviewID, err := uuid.Parse(req.ID)
	if err != nil {
		err2 := apierror.ErrValidation
		apierror.AddPayload(&err2, "INVALID_UUID")
		return err2
	}

	review, err := uc.repo.GetByID(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReviewNotFound
		}
		log.Print("Error getting review by id: ", err)
		return apierror.ErrInternalServer
	}

	if review.UserID != userID {
		return apierror.ErrNotYourResource
	}

	rating, count, err := uc.courseRepo.GetRating(ctx, review.CourseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return course.ErrCourseNotFound
		}
		log.Print("Error getting course rating: ", err)
		return apierror.ErrInternalServer
	}

	// mean_{updated} = mean_{n} + \frac{x_{new} - x_{old}}{n}
	var newRating float32
	if req.Rating != 0 {
		newRating = rating + (float32(req.Rating)-float32(review.Rating))/float32(count)
	}

	newReview := &schema.Review{
		ID:       reviewID,
		Rating:   req.Rating,
		Feedback: req.Feedback,
	}

	if err = uc.repo.Update(newReview, review.CourseID, newRating); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReviewNotFound
		}
		log.Print("Error updating review: ", err)
		return apierror.ErrInternalServer
	}

	return nil
}

func (uc *UseCase) Delete(ctx context.Context, req *DeleteReviewRequest) error {
	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return apierror.ErrInternalServer
	}

	reviewID, err := uuid.Parse(req.ID)
	if err != nil {
		err2 := apierror.ErrValidation
		apierror.AddPayload(&err2, "INVALID_UUID")
		return err2
	}

	review, err := uc.repo.GetByID(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReviewNotFound
		}
		log.Print("Error getting review by id: ", err)
		return apierror.ErrInternalServer
	}

	if review.UserID != userID {
		return apierror.ErrNotYourResource
	}

	err = uc.repo.Delete(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReviewNotFound
		}
		log.Print("Error deleting review: ", err)
		return apierror.ErrInternalServer
	}

	return nil
}
