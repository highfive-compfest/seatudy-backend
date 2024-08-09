package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/config"
	"github.com/highfive-compfest/seatudy-backend/internal/fileutil"
	"gorm.io/gorm"
	"log"
	"slices"
	"strings"
	"time"
)

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetByID(req *GetUserByIDRequest) (*GetUserResponse, error) {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		err2 := apierror.ErrValidation
		apierror.AddPayload(&err2, "INVALID_UUID")
		return nil, err2
	}

	user, err := uc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		log.Println("Error getting user by id: ", err)
		return nil, apierror.ErrInternalServer
	}

	return &GetUserResponse{
		ID:              user.ID.String(),
		Email:           user.Email,
		Name:            user.Name,
		ImageURL:        user.ImageURL,
		Role:            string(user.Role),
		IsEmailVerified: user.IsEmailVerified,
		CreatedAt:       user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (uc *UseCase) Update(ctx context.Context, req *UpdateUserRequest) error {
	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		log.Println("Error parsing uuid from jwt: ", err)
		return apierror.ErrInternalServer
	}

	var imageUrl string
	if req.ImageFile != nil && req.ImageFile.Size > 0 {
		if req.ImageFile.Size > 2*fileutil.MegaByte {
			err2 := apierror.ErrFileTooLarge
			apierror.AddPayload(&err2, map[string]string{
				"max_size":      "2 MB",
				"received_size": fileutil.ByteToAppropriateUnit(req.ImageFile.Size),
			})
			return err2
		}

		fileType, err := fileutil.DetectMultipartFileType(req.ImageFile)
		if err != nil {
			log.Println("Error detecting image type: ", err)
			return apierror.ErrInternalServer
		}
		allowedTypes := fileutil.ImageContentTypes
		if !slices.Contains(allowedTypes, fileType) {
			err2 := apierror.ErrInvalidFileType
			apierror.AddPayload(&err2, map[string]any{
				"allowed_types": allowedTypes,
				"received_type": fileType,
			})
			return err2
		}

		imageUrl, err = config.UploadFile(
			"users/avatar/"+userID.String()+"."+strings.Split(fileType, "/")[1],
			req.ImageFile,
		)

		if err != nil {
			log.Println("Error uploading image: ", err)
			return apierror.ErrInternalServer
		}
	}

	userEntity := User{
		ID:       userID,
		Name:     req.Name,
		ImageURL: imageUrl,
	}

	if err := uc.repo.Update(&userEntity); err != nil {
		log.Println("Error updating user: ", err)
		return apierror.ErrInternalServer
	}

	return nil
}
