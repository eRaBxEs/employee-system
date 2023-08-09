package storage

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"employee-management-system/model"
	"employee-management-system/pkg/helper"
)

// UserDatabase enlist all possible storage operations for Users
//
//go:generate mockgen -source user.go -destination ./mock/mock_user.go -package mock UserDatabase
type UserDatabase interface {
	Register(ctx context.Context, user model.User) (model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error)
	Authenticate(ctx context.Context, email, password string) (*model.User, error)
}

// User object
type User struct {
	logger  zerolog.Logger
	storage *Storage
}

// NewUser creates a new reference to the User storage entity
func NewUser(s *Storage) *UserDatabase {
	l := s.Logger.With().Str(helper.LogStrKeyLevel, "user").Logger()
	user := &User{
		logger:  l,
		storage: s,
	}
	userDB := UserDatabase(user)
	return &userDB
}

// Register or create a new user into the storage
func (u *User) Register(ctx context.Context, user model.User) (model.User, error) {

	db := u.storage.DB.WithContext(ctx).Create(&user)
	if db.Error != nil {
		u.logger.Err(db.Error).Msgf("User::Register error: %v, (%v)", ErrRecordCreatingFailed, db.Error)
		if strings.Contains(db.Error.Error(), "duplicate key value") {
			return model.User{}, ErrDuplicateRecord
		}
		return model.User{}, ErrRecordCreatingFailed
	}
	return user, nil
}

// GetUserByID should find a user by it's ID
func (u *User) GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	var user model.User
	db := u.storage.DB.WithContext(ctx).Where("id = ?", id.String()).Find(&user)
	if db.Error != nil {
		u.logger.Err(db.Error).Msgf("User::GetUserByID error: %v, (%v)", ErrRecordNotFound, db.Error)
		return user, ErrRecordNotFound
	}
	return user, nil
}

// Authenticate tests supplied username and password to attempt login against the user table
func (u *User) Authenticate(ctx context.Context, email, password string) (*model.User, error) {
	var user model.User
	db := u.storage.DB.WithContext(ctx).Where("user_name = ?", email).Find(&user)
	if db.Error != nil {
		u.logger.Err(db.Error).Msgf("User::Authenticate error: %v, (%v)", ErrRecordNotFound, db.Error)
		return nil, ErrRecordNotFound
	}

	// if user is found using the email address
	if user.Password.Check(model.Password(password)) {
		return &user, nil
	}

	return nil, ErrPasswordIncorrect
}
