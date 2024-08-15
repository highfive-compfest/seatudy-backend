package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/config"
	"github.com/highfive-compfest/seatudy-backend/internal/fileutil"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"gorm.io/gorm"
	"log"
	"slices"
	"strings"
)

type UseCase struct {
	repo     IRepository
	uploader config.FileUploader
}

func NewUseCase(repo IRepository, uploader config.FileUploader) *UseCase {
	return &UseCase{repo: repo, uploader: uploader}
}

func (uc *UseCase) GetMe(ctx context.Context) (*schema.User, error) {
	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		return nil, apierror.ErrTokenInvalid.Build()
	}

	user, err := uc.repo.GetByID(userID)
	if err != nil {
		log.Println("Error getting user by id: ", err)
		return nil, apierror.ErrInternalServer.Build()
	}

	return user, nil
}

func (uc *UseCase) GetByID(req *GetUserByIDRequest) (*GetUserResponse, error) {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, apierror.ErrInvalidParamId.Build()
	}

	user, err := uc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound.Build()
		}
		log.Println("Error getting user by id: ", err)
		return nil, apierror.ErrInternalServer.Build()
	}

	return &GetUserResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		ImageURL: user.ImageURL,
		Role:     string(user.Role),
	}, nil
}

func (uc *UseCase) Update(ctx context.Context, req *UpdateUserRequest) error {
	userID, err := uuid.Parse(ctx.Value("user.id").(string))
	if err != nil {
		log.Println("Error parsing uuid from jwt: ", err)
		return apierror.ErrTokenInvalid.Build()
	}

	var imageUrl string
	if req.ImageFile != nil && req.ImageFile.Size > 0 {
		if req.ImageFile.Size > 2*fileutil.MegaByte {
			err2 := apierror.ErrFileTooLarge.WithPayload(map[string]string{
				"max_size":      "2 MB",
				"received_size": fileutil.ByteToAppropriateUnit(req.ImageFile.Size),
			})
			return err2.Build()
		}

		fileType, err := fileutil.DetectMultipartFileType(req.ImageFile)
		if err != nil {
			log.Println("Error detecting image type: ", err)
			return apierror.ErrInvalidFileType.Build()
		}
		allowedTypes := fileutil.ImageContentTypes
		if !slices.Contains(allowedTypes, fileType) {
			err2 := apierror.ErrInvalidFileType.WithPayload(map[string]any{
				"allowed_types": allowedTypes,
				"received_type": fileType,
			})
			return err2.Build()
		}

		imageUrl, err = uc.uploader.UploadFile(
			"users/avatar/"+userID.String()+"."+strings.Split(fileType, "/")[1],
			req.ImageFile,
		)

		if err != nil {
			log.Println("Error uploading image: ", err)
			return apierror.ErrInternalServer.Build()
		}
	}

	userEntity := schema.User{
		ID:       userID,
		Name:     req.Name,
		ImageURL: imageUrl,
	}

	if err := uc.repo.Update(&userEntity); err != nil {
		log.Println("Error updating user: ", err)
		return apierror.ErrInternalServer.Build()
	}

	return nil
}
