package forum

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/course"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/courseenroll"
	"github.com/highfive-compfest/seatudy-backend/internal/pagination"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"gorm.io/gorm"
	"log"
)

type UseCase struct {
	repo       IRepository
	enrollUc   *courseenroll.UseCase
	courseRepo course.Repository
}

func NewUseCase(repo IRepository, enrollUc *courseenroll.UseCase, courseRepo course.Repository) *UseCase {
	return &UseCase{repo: repo, enrollUc: enrollUc, courseRepo: courseRepo}
}

func (uc *UseCase) isPermitted(ctx context.Context, userRole string, userID, courseID uuid.UUID) (bool, error) {
	if userRole == "instructor" {
		courseObj, err := uc.courseRepo.GetByID(ctx, courseID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, course.ErrCourseNotFound
			}
			log.Println("Error getting course: ", err)
			return false, apierror.ErrInternalServer
		}

		if courseObj.InstructorID != userID {
			return false, apierror.ErrForbidden
		}
	} else {
		ok, err := uc.enrollUc.CheckEnrollment(ctx, userID, courseID)
		if err != nil {
			log.Println("Error checking enrollment: ", err)
			return false, apierror.ErrInternalServer
		}
		if !ok {
			return false, courseenroll.ErrNotEnrolled
		}
	}

	return true, nil
}

func (uc *UseCase) CreateDiscussion(ctx context.Context, req *CreateForumDiscussionRequest) error {
	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return apierror.ErrTokenInvalid
	}

	userRole := ctx.Value("user.role").(string)

	ok, err := uc.isPermitted(ctx, userRole, userID, req.CourseID)
	if err != nil {
		return err
	}
	if !ok {
		return courseenroll.ErrNotEnrolled
	}

	discussionID, err := uuid.NewV7()
	if err != nil {
		return apierror.ErrInternalServer
	}

	discussion := &schema.ForumDiscussion{
		ID:       discussionID,
		UserID:   userID,
		CourseID: req.CourseID,
		Title:    req.Title,
		Content:  req.Content,
	}

	if err := uc.repo.CreateDiscussion(discussion); err != nil {
		log.Println("Error creating discussion: ", err)
		return apierror.ErrInternalServer
	}

	return nil
}

func (uc *UseCase) GetDiscussionByID(ctx context.Context, idStr string) (*schema.ForumDiscussion, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, apierror.ErrValidation
	}

	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return nil, apierror.ErrTokenInvalid
	}

	discussion, err := uc.repo.GetDiscussionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDiscussionNotFound
		}
		log.Println("Error getting discussion: ", err)
		return nil, apierror.ErrInternalServer
	}

	ok, err := uc.isPermitted(ctx, ctx.Value("user.role").(string), userID, discussion.CourseID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, courseenroll.ErrNotEnrolled
	}

	return discussion, nil
}

func (uc *UseCase) GetDiscussionsByCourseID(ctx context.Context, req *GetForumDiscussionsRequest) (*pagination.GetResourcePaginatedResponse, error) {
	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return nil, apierror.ErrTokenInvalid
	}

	courseID, err := uuid.Parse(req.CourseID)
	if err != nil {
		return nil, apierror.ErrValidation
	}

	ok, err := uc.isPermitted(ctx, ctx.Value("user.role").(string), userID, courseID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, courseenroll.ErrNotEnrolled
	}

	discussions, total, err := uc.repo.GetDiscussionsByCourseID(courseID, req.Page, req.Limit)
	if err != nil {
		log.Println("Error getting discussions: ", err)
		return nil, apierror.ErrInternalServer
	}

	resp := pagination.GetResourcePaginatedResponse{
		Data:       discussions,
		Pagination: pagination.NewPagination(int(total), req.Page, req.Limit),
	}

	return &resp, nil
}

func (uc *UseCase) UpdateDiscussion(ctx context.Context, req *UpdateForumDiscussionRequest) error {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return apierror.ErrValidation
	}

	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return apierror.ErrValidation
	}

	discussion, err := uc.repo.GetDiscussionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDiscussionNotFound
		}
		log.Println("Error getting discussion: ", err)
		return apierror.ErrInternalServer
	}

	if discussion.UserID != userID {
		return apierror.ErrNotYourResource
	}

	newDiscussion := &schema.ForumDiscussion{
		ID:      id,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := uc.repo.UpdateDiscussion(newDiscussion); err != nil {
		log.Println("Error updating discussion: ", err)
		return apierror.ErrInternalServer
	}

	return nil
}

func (uc *UseCase) DeleteDiscussion(ctx context.Context, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return apierror.ErrValidation
	}

	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return apierror.ErrValidation
	}

	discussion, err := uc.repo.GetDiscussionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDiscussionNotFound
		}
		log.Println("Error getting discussion: ", err)
		return apierror.ErrInternalServer
	}

	if discussion.UserID != userID {
		return apierror.ErrNotYourResource
	}

	if err := uc.repo.DeleteDiscussion(id); err != nil {
		log.Println("Error deleting discussion: ", err)
		return apierror.ErrInternalServer
	}

	return nil
}

func (uc *UseCase) CreateReply(ctx context.Context, req *CreateForumReplyRequest) error {
	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return apierror.ErrTokenInvalid
	}

	discussion, err := uc.repo.GetDiscussionByID(req.DiscussionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDiscussionNotFound
		}
		log.Println("Error getting discussion: ", err)
		return apierror.ErrInternalServer
	}

	ok, err := uc.isPermitted(ctx, ctx.Value("user.role").(string), userID, discussion.CourseID)
	if err != nil {
		return err
	}
	if !ok {
		return courseenroll.ErrNotEnrolled
	}

	replyID, err := uuid.NewV7()
	if err != nil {
		return apierror.ErrInternalServer
	}

	reply := &schema.ForumReply{
		ID:                replyID,
		UserID:            userID,
		ForumDiscussionID: discussion.ID,
		CourseID:          discussion.CourseID,
		Content:           req.Content,
	}

	if err := uc.repo.CreateReply(reply); err != nil {
		log.Println("Error creating reply: ", err)
		return apierror.ErrInternalServer
	}

	return nil
}

func (uc *UseCase) GetReplyByID(ctx context.Context, idStr string) (*schema.ForumReply, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, apierror.ErrValidation
	}

	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return nil, apierror.ErrTokenInvalid
	}

	reply, err := uc.repo.GetReplyByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReplyNotFound
		}
		log.Println("Error getting reply: ", err)
		return nil, apierror.ErrInternalServer
	}

	ok, err := uc.isPermitted(ctx, ctx.Value("user.role").(string), userID, reply.CourseID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, courseenroll.ErrNotEnrolled
	}

	return reply, nil
}

func (uc *UseCase) GetRepliesByDiscussionID(ctx context.Context, req *GetForumRepliesRequest) (*pagination.GetResourcePaginatedResponse, error) {
	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return nil, apierror.ErrTokenInvalid
	}

	discussionID, err := uuid.Parse(req.DiscussionID)
	if err != nil {
		return nil, apierror.ErrValidation
	}

	discussion, err := uc.repo.GetDiscussionByID(discussionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDiscussionNotFound
		}
		log.Println("Error getting discussion: ", err)
		return nil, apierror.ErrInternalServer
	}

	ok, err := uc.isPermitted(ctx, ctx.Value("user.role").(string), userID, discussion.CourseID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, courseenroll.ErrNotEnrolled
	}

	replies, total, err := uc.repo.GetRepliesByDiscussionID(discussionID, req.Page, req.Limit)
	if err != nil {
		log.Println("Error getting replies: ", err)
		return nil, apierror.ErrInternalServer
	}

	resp := pagination.GetResourcePaginatedResponse{
		Data:       replies,
		Pagination: pagination.NewPagination(int(total), req.Page, req.Limit),
	}

	return &resp, nil
}

func (uc *UseCase) UpdateReply(ctx context.Context, req *UpdateForumReplyRequest) error {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return apierror.ErrValidation
	}

	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return apierror.ErrValidation
	}

	reply, err := uc.repo.GetReplyByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReplyNotFound
		}
		log.Println("Error getting reply: ", err)
		return apierror.ErrInternalServer
	}

	if reply.UserID != userID {
		return apierror.ErrNotYourResource
	}

	newReply := &schema.ForumReply{
		ID:      id,
		Content: req.Content,
	}

	if err := uc.repo.UpdateReply(newReply); err != nil {
		log.Println("Error updating reply: ", err)
		return apierror.ErrInternalServer
	}

	return nil
}

func (uc *UseCase) DeleteReply(ctx context.Context, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return apierror.ErrValidation
	}

	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return apierror.ErrTokenInvalid
	}

	reply, err := uc.repo.GetReplyByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReplyNotFound
		}
		log.Println("Error getting reply: ", err)
		return apierror.ErrInternalServer
	}

	if reply.UserID != userID {
		return apierror.ErrNotYourResource
	}

	if err := uc.repo.DeleteReply(id); err != nil {
		log.Println("Error deleting reply: ", err)
		return apierror.ErrInternalServer
	}

	return nil
}
