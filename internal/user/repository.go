package user

import "gorm.io/gorm"

type Repository interface {
	// CreateUser creates a new user.
	CreateUser(user *User) error

	// GetUserByEmail retrieves a user by their email.
	GetUserByEmail(email string) (*User, error)

	// UpdateUser updates user information.
	UpdateUser(user *User) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateUser(user *User) error {
	// Implementation here
	return nil
}

func (r *repository) GetUserByEmail(email string) (*User, error) {
	// Implementation here
	return nil, nil
}

func (r *repository) UpdateUser(user *User) error {
	// Implementation here
	return nil
}
