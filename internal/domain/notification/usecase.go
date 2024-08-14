package notification

import (
	"context"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/pagination"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"log"
)

type UseCase struct {
	repo IRepository
}

func NewUseCase(repo IRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Create(notification *schema.Notification) error {
	if err := uc.repo.Create(notification); err != nil {
		log.Println("Error creating notification: ", err)
		return apierror.ErrInternalServer
	}

	return nil
}

func (uc *UseCase) GetMy(ctx context.Context, req *GetByUserIDRequest) (*pagination.GetResourcePaginatedResponse, error) {
	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return nil, apierror.ErrTokenInvalid
	}

	offset := (req.Page - 1) * req.Limit
	notifications, total, err := uc.repo.GetByUserID(userID, req.Limit, offset)
	if err != nil {
		log.Println("Error getting notifications: ", err)
		return nil, apierror.ErrInternalServer
	}

	res := &pagination.GetResourcePaginatedResponse{
		Data:       notifications,
		Pagination: pagination.NewPagination(int(total), req.Page, req.Limit),
	}

	return res, nil
}

func (uc *UseCase) GetUnreadCount(ctx context.Context) (int64, error) {
	id, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return 0, apierror.ErrTokenInvalid
	}

	count, err := uc.repo.GetUnreadCount(id)
	if err != nil {
		log.Println("Error getting unread count: ", err)
		return 0, apierror.ErrInternalServer
	}

	return count, nil
}

func (uc *UseCase) UpdateRead(notificationID string) error {
	id, err := uuid.Parse(notificationID)
	if err != nil {
		return apierror.ErrInvalidParamId
	}

	if err := uc.repo.UpdateRead(id); err != nil {
		log.Println("Error updating notification: ", err)
		return apierror.ErrInternalServer
	}

	return nil
}