package auth

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Repository interface {
	SaveOTP(ctx context.Context, email string, otp string) error
	GetOTP(ctx context.Context, email string) (string, error)
	DeleteOTP(ctx context.Context, email string) error
	SaveResetPasswordToken(ctx context.Context, email string, token string) error
	GetResetPasswordToken(ctx context.Context, email string) (string, error)
	DeleteResetPasswordToken(ctx context.Context, email string) error
}

type repository struct {
	rds *redis.Client
}

func NewRepository(rds *redis.Client) Repository {
	return &repository{rds: rds}
}

func (r *repository) SaveOTP(ctx context.Context, email string, otp string) error {
	return r.rds.Set(ctx, "auth:"+email+":otp", otp, 10*time.Minute).Err()
}

func (r *repository) GetOTP(ctx context.Context, email string) (string, error) {
	return r.rds.Get(ctx, "auth:"+email+":otp").Result()
}

func (r *repository) DeleteOTP(ctx context.Context, email string) error {
	return r.rds.Del(ctx, "auth:"+email+":otp").Err()
}

func (r *repository) SaveResetPasswordToken(ctx context.Context, email string, token string) error {
	return r.rds.Set(ctx, "auth:"+email+":reset_password_token", token, 10*time.Minute).Err()
}

func (r *repository) GetResetPasswordToken(ctx context.Context, email string) (string, error) {
	return r.rds.Get(ctx, "auth:"+email+":reset_password_token").Result()
}

func (r *repository) DeleteResetPasswordToken(ctx context.Context, email string) error {
	return r.rds.Del(ctx, "auth:"+email+":reset_password_token").Err()
}
