package auth

import "gorm.io/gorm"

type Repository interface {
	// StoreOTP stores the OTP in Redis.
	StoreOTP(email, otp string) error

	// GetOTP retrieves the OTP from Redis.
	GetOTP(email string) (string, error)

	// To be continued
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) StoreOTP(email, otp string) error {
	// Implementation here
	return nil
}

func (r *repository) GetOTP(email string) (string, error) {
	// Implementation here
	return "", nil
}
