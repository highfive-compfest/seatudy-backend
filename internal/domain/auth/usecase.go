package auth

import (
	"context"
	_ "embed"
	"errors"
	"github.com/google/uuid"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/config"
	"github.com/highfive-compfest/seatudy-backend/internal/domain/user"
	"github.com/highfive-compfest/seatudy-backend/internal/jwtoken"
	"github.com/highfive-compfest/seatudy-backend/internal/mailer"
	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

type UseCase struct {
	authRepo   Repository
	userRepo   user.IRepository
	mailDialer config.IMailer
}

func NewUseCase(authRepo Repository, userRepo user.IRepository, mailDialer config.IMailer) *UseCase {
	return &UseCase{authRepo: authRepo, userRepo: userRepo, mailDialer: mailDialer}
}

func (uc *UseCase) Register(req *RegisterRequest) error {
	// Hash & Salt password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password: ", err)
		return apierror.ErrInternalServer.Build()
	}

	id, err := uuid.NewV7()
	if err != nil {
		log.Println("Error generating UUID: ", err)
		return apierror.ErrInternalServer.Build()
	}

	userEntity := schema.User{
		ID:              id,
		Email:           req.Email,
		IsEmailVerified: false,
		Name:            req.Name,
		PasswordHash:    string(passwordHash),
		Role:            schema.Role(req.Role),
		ImageURL:        "",
	}

	err = uc.userRepo.Create(&userEntity)
	if err != nil {
		var pgErr *pgconn.PgError
		ok := errors.As(err, &pgErr)
		if ok && pgErr.Code == "23505" {
			return ErrEmailAlreadyRegistered.Build()
		}
		log.Println("Error creating user: ", err)
		return apierror.ErrInternalServer.Build()
	}

	return nil
}

func (uc *UseCase) Login(req *LoginRequest) (*LoginResponse, error) {
	usr, err := uc.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials.Build()
		}
		log.Println("Error getting user by email: ", err)
		return nil, apierror.ErrInternalServer.Build()
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials.Build()
	}

	accessToken, err := jwtoken.CreateAccessJWT(
		usr.ID.String(), usr.Email, usr.IsEmailVerified, usr.Name, string(usr.Role),
	)
	if err != nil {
		log.Println("Error creating access token: ", err)
		return nil, apierror.ErrInternalServer.Build()
	}

	refreshToken, err := jwtoken.CreateRefreshJWT(usr.ID.String())
	if err != nil {
		log.Println("Error creating refresh token: ", err)
		return nil, apierror.ErrInternalServer.Build()
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         usr,
	}, nil
}

func (uc *UseCase) Refresh(req *RefreshRequest) (*RefreshResponse, error) {
	claims, err := jwtoken.DecodeRefreshJWT(req.RefreshToken)
	if err != nil {
		return nil, apierror.ErrTokenInvalid.Build()
	}

	if claims.Issuer != "seatudy-backend-refreshtoken" {
		return nil, apierror.ErrTokenInvalid.Build()
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, apierror.ErrTokenExpired.Build()
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, apierror.ErrTokenInvalid.Build()
	}

	userEntity, err := uc.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.ErrTokenInvalid.Build()
		}
		log.Println("Error getting user by ID: ", err)
		return nil, apierror.ErrInternalServer.Build()
	}

	accessToken, err := jwtoken.CreateAccessJWT(
		userEntity.ID.String(), userEntity.Email, userEntity.IsEmailVerified, userEntity.Name, string(userEntity.Role),
	)
	if err != nil {
		log.Println("Error creating access token: ", err)
		return nil, apierror.ErrInternalServer.Build()
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

//go:embed otp_email_template.html
var otpEmailTemplate string

func (uc *UseCase) SendOTP(ctx context.Context) error {
	email := ctx.Value("user.email").(string)
	name := ctx.Value("user.name").(string)
	isEmailVerified := ctx.Value("user.is_email_verified").(bool)

	if isEmailVerified {
		return ErrEmailAlreadyVerified.Build()
	}

	otp := strconv.Itoa(generateOTP())

	// Save OTP to database
	err := uc.authRepo.SaveOTP(ctx, email, otp)
	if err != nil {
		log.Println("Error saving OTP: ", err)
		return apierror.ErrInternalServer.Build()
	}

	// Send OTP to email
	data := map[string]any{
		"recipient_name": name,
		"otp":            otp,
	}

	mail, err := mailer.GenerateMail(email, "Your Seatudy OTP Code", otpEmailTemplate, data)
	if err != nil {
		log.Println("Error generating OTP email: ", err)
		return apierror.ErrInternalServer.Build()
	}

	if err = uc.mailDialer.DialAndSend(mail); err != nil {
		log.Println("Error sending OTP email: ", err)
		return apierror.ErrInternalServer.Build()
	}

	return nil
}

func (uc *UseCase) VerifyOTP(ctx context.Context, req *VerifyEmailRequest) error {
	email := ctx.Value("user.email").(string)

	savedOTP, err := uc.authRepo.GetOTP(ctx, email)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrExpiredOTP.Build()
		}
		log.Println("Error getting OTP: ", err)
		return apierror.ErrInternalServer.Build()
	}

	if req.OTP != savedOTP {
		return ErrInvalidOTP.Build()
	}

	if err = uc.authRepo.DeleteOTP(ctx, email); err != nil {
		log.Println("Error deleting OTP: ", err)
		return apierror.ErrInternalServer.Build()
	}

	if err := uc.userRepo.UpdateByEmail(email, &schema.User{IsEmailVerified: true}); err != nil {
		log.Println("Error updating user email verification status: ", err)
		return apierror.ErrInternalServer.Build()
	}

	return nil
}

// generateRandomString generates a url safe random string of length n
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

//go:embed reset_password_email_template.html
var resetPasswordEmailTemplate string

func (uc *UseCase) SendResetPasswordLink(ctx context.Context, req *SendResetPasswordLinkRequest) error {
	userEntity, err := uc.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.ErrUserNotFound.Build()
		}
		log.Println("Error getting user by email: ", err)
		return apierror.ErrInternalServer.Build()
	}

	token := generateRandomString(32)

	err = uc.authRepo.SaveResetPasswordToken(ctx, req.Email, token)
	if err != nil {
		log.Println("Error saving reset password token: ", err)
		return apierror.ErrInternalServer.Build()
	}

	data := map[string]any{
		"recipient_name": userEntity.Name,
		"reset_link": config.Env.FrontendUrl +
			"/reset-password?token=" + token + "&email=" + url.QueryEscape(req.Email),
	}

	mail, err := mailer.GenerateMail(req.Email, "Reset Your Seatudy Password", resetPasswordEmailTemplate, data)
	if err != nil {
		log.Println("Error generating reset password email: ", err)
		return apierror.ErrInternalServer.Build()
	}

	if err := uc.mailDialer.DialAndSend(mail); err != nil {
		log.Println("Error sending reset password email: ", err)
		return apierror.ErrInternalServer.Build()
	}

	return nil
}

func (uc *UseCase) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	savedToken, err := uc.authRepo.GetResetPasswordToken(ctx, req.Email)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrExpiredResetPasswordLink.Build()
		}
		log.Println("Error getting reset password token: ", err)
		return apierror.ErrInternalServer.Build()
	}

	if req.Token != savedToken {
		return ErrInvalidResetPasswordLink.Build()
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password: ", err)
		return apierror.ErrInternalServer.Build()
	}

	err = uc.userRepo.UpdateByEmail(req.Email, &schema.User{PasswordHash: string(passwordHash)})
	if err != nil {
		log.Println("Error updating user password: ", err)
		return apierror.ErrInternalServer.Build()
	}

	if err = uc.authRepo.DeleteResetPasswordToken(ctx, req.Email); err != nil {
		log.Println("Error deleting reset password token: ", err)
		return apierror.ErrInternalServer.Build()
	}

	return nil
}

func (uc *UseCase) ChangePassword(ctx context.Context, req *ChangePasswordRequest) error {
	email := ctx.Value("user.email").(string)

	userEntity, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.ErrUserNotFound.Build()
		}
		log.Println("Error getting user by email: ", err)
		return apierror.ErrInternalServer.Build()
	}

	err = bcrypt.CompareHashAndPassword([]byte(userEntity.PasswordHash), []byte(req.OldPassword))
	if err != nil {
		return ErrInvalidCredentials.Build()
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password: ", err)
		return apierror.ErrInternalServer.Build()
	}

	err = uc.userRepo.UpdateByEmail(email, &schema.User{PasswordHash: string(passwordHash)})
	if err != nil {
		log.Println("Error updating user password: ", err)
		return apierror.ErrInternalServer.Build()
	}

	return nil
}
