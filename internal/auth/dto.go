package auth

import (
	"github.com/highfive-compfest/seatudy-backend/internal/user"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=320"`
	Name     string `json:"full_name" binding:"required,max=50"`
	Password string `json:"password" binding:"required,max=72,min=8"`
	Role     string `json:"role" binding:"required,oneof=student instructor"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=320"`
	Password string `json:"password" binding:"required,max=72"`
}

type RegisterLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	user.User    `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type VerifyEmailRequest struct {
	OTP string `json:"otp" binding:"required"`
}

type VerifyEmailResponse struct {
	AccessToken string `json:"access_token"`
}
