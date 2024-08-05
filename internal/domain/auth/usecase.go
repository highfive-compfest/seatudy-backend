package auth

import (
	"bytes"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/user"
	"github.com/highfive-compfest/seatudy-backend/internal/jwtoken"
	"github.com/highfive-compfest/seatudy-backend/internal/mailer"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"html/template"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type UseCase struct {
	authRepo   Repository
	userRepo   user.Repository
	mailDialer *gomail.Dialer
}

func NewUseCase(authRepo Repository, userRepo user.Repository, mailDialer *gomail.Dialer) *UseCase {
	return &UseCase{authRepo: authRepo, userRepo: userRepo, mailDialer: mailDialer}
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
		return nil, apierror.ErrTokenInvalid
	}

	if claims.Issuer != "seatudy-backend-refreshtoken" {
		return nil, apierror.ErrTokenInvalid
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, apierror.ErrTokenExpired
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, apierror.ErrTokenInvalid
	}

	userEntity, err := uc.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.ErrTokenInvalid
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

func generateOTP() int {
	high := 999999
	low := 100000
	return rand.Intn(high-low) + low
}

func generateOTPMail(email, name, otp string) (*gomail.Message, error) {
	data := map[string]any{
		"recipient_name": name,
		"otp":            otp,
	}

	tmpl, err := template.ParseFiles("internal/domain/auth/otp_email_template.html")
	if err != nil {
		return nil, err
	}

	var tmplOutput bytes.Buffer
	err = tmpl.Execute(&tmplOutput, data)
	if err != nil {
		return nil, err
	}

	mail := mailer.NewMail()
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "Your Seatudy OTP Code")
	mail.SetBody("text/html", tmplOutput.String())

	return mail, nil
}

func (uc *UseCase) SendOTP(ctx context.Context) error {
	email := ctx.Value("user.email").(string)
	name := ctx.Value("user.name").(string)
	isEmailVerified := ctx.Value("user.is_email_verified").(bool)

	if isEmailVerified {
		return ErrEmailAlreadyVerified
	}

	otp := strconv.Itoa(generateOTP())

	// Save OTP to database
	err := uc.authRepo.SaveOTP(ctx, email, otp)
	if err != nil {
		log.Println("Error saving OTP: ", err)
		return apierror.ErrInternalServer
	}

	// Send OTP to email
	mail, err := generateOTPMail(email, name, otp)
	if err != nil {
		log.Println("Error generating OTP email: ", err)
		return apierror.ErrInternalServer
	}

	err = uc.mailDialer.DialAndSend(mail)
	if err != nil {
		log.Println("Error sending OTP email: ", err)
		return apierror.ErrInternalServer
	}

	return nil
}

func (uc *UseCase) VerifyOTP(ctx context.Context, otp string) error {
	email := ctx.Value("user.email").(string)

	savedOTP, err := uc.authRepo.GetOTP(ctx, email)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrExpiredOTP
		}
		log.Println("Error getting OTP: ", err)
		return apierror.ErrInternalServer
	}

	if otp != savedOTP {
		return ErrInvalidOTP
	}

	if err = uc.authRepo.DeleteOTP(ctx, email); err != nil {
		log.Println("Error deleting OTP: ", err)
		return apierror.ErrInternalServer
	}

	if err := uc.userRepo.UpdateByEmail(email, &user.User{IsEmailVerified: true}); err != nil {
		log.Println("Error updating user email verification status: ", err)
		return apierror.ErrInternalServer
	}

	return nil
}
