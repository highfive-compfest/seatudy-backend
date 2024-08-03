package auth

import (
	"context"
	"github.com/highfive-compfest/seatudy-backend/internal/user"
)

type UseCase struct {
	authRepo Repository
	userRepo user.Repository
}

func NewUseCase(authRepo Repository, userRepo user.Repository) *UseCase {
	return &UseCase{authRepo: authRepo}
}

func (uc *UseCase) Register(ctx context.Context, req RegisterRequest) (*RegisterLoginResponse, error) {
	// Implementation here
	return nil, nil
}

func (uc *UseCase) Login(ctx context.Context, req LoginRequest) (*RegisterLoginResponse, error) {
	// Implementation here
	return nil, nil
}

func (uc *UseCase) Refresh(ctx context.Context, req RefreshRequest) (*RefreshResponse, error) {
	// Implementation here
	return nil, nil
}

func (uc *UseCase) VerifyEmail(ctx context.Context, req VerifyEmailRequest) (*VerifyEmailResponse, error) {
	// Implementation here
	return nil, nil
}
