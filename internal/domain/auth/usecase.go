package auth

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/user"
	"github.com/highfive-compfest/seatudy-backend/internal/jwtoken"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"time"
)

type UseCase struct {
	authRepo Repository
	userRepo user.Repository
}

func NewUseCase(authRepo Repository, userRepo user.Repository) *UseCase {
	return &UseCase{authRepo: authRepo, userRepo: userRepo}
}

func (uc *UseCase) Register(req RegisterRequest) error {
	// Hash & Salt password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password: ", err)
		return apierror.ErrInternalServer
	}

	id, err := uuid.NewV7()
	if err != nil {
		log.Println("Error generating UUID: ", err)
		return apierror.ErrInternalServer
	}

	userEntity := user.User{
		ID:              id,
		Email:           req.Email,
		IsEmailVerified: false,
		Name:            req.Name,
		PasswordHash:    string(passwordHash),
		Role:            user.Role(req.Role),
		ImageURL:        "",
		Balance:         0,
	}

	err = uc.userRepo.Create(&userEntity)
	if err != nil {
		var pgErr *pgconn.PgError
		ok := errors.As(err, &pgErr)
		if ok && pgErr.Code == "23505" {
			return ErrEmailAlreadyRegistered
		}
		log.Println("Error creating user: ", err)
		return apierror.ErrInternalServer
	}

	return nil
}

func (uc *UseCase) Login(req LoginRequest) (*LoginResponse, error) {
	userEntity, err := uc.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		log.Println("Error getting user by email: ", err)
		return nil, apierror.ErrInternalServer
	}

	err = bcrypt.CompareHashAndPassword([]byte(userEntity.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := jwtoken.CreateAccessJWT(userEntity)
	if err != nil {
		log.Println("Error creating access token: ", err)
		return nil, apierror.ErrInternalServer
	}

	refreshToken, err := jwtoken.CreateRefreshJWT(userEntity)
	if err != nil {
		log.Println("Error creating refresh token: ", err)
		return nil, apierror.ErrInternalServer
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userEntity,
	}, nil
}

func (uc *UseCase) Refresh(req RefreshRequest) (*RefreshResponse, error) {
	claims, err := jwtoken.DecodeRefreshJWT(req.RefreshToken)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	if claims.Issuer != "seatudy-backend-refreshtoken" {
		return nil, ErrTokenInvalid
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	userEntity, err := uc.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenInvalid
		}
		log.Println("Error getting user by ID: ", err)
		return nil, apierror.ErrInternalServer
	}

	accessToken, err := jwtoken.CreateAccessJWT(userEntity)
	if err != nil {
		log.Println("Error creating access token: ", err)
		return nil, apierror.ErrInternalServer
	}

	return &RefreshResponse{
		AccessToken: accessToken,
	}, nil
}

func (uc *UseCase) VerifyEmail(ctx context.Context, req VerifyEmailRequest) (*VerifyEmailResponse, error) {
	// Implementation here
	return nil, nil
}
