package model

import (
	"time"

	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

// Kind is an integer value of supported kinds of users on this platform
type Kind int

// Password is string representation of user password either when it is encrypted or not
type Password string

const (
	// KindUnknown is an invalid or unknown kind of user
	KindUnknown Kind = 0
	// KindAdministrator is an administrative kind of user
	KindAdministrator = 1
	// KindStaff is a staff kind of user
	KindStaff = 2
)

// User object
type (
	User struct {
		ID        int `gorm:"column:id;PRIMARY_KEY;type:int;"`
		UserName  *string
		Password  Password
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *gorm.DeletedAt
	}
)

// String representation of a user Kind int value
func (k Kind) String() string {
	return [...]string{
		"Unknown",
		"Administrator",
		"Agent",
		"Partner",
	}[k]
}

// Value Get the int value of type Kind
func (k Kind) Value() int {
	return int(k)
}

// String representation of a password either as encrypted or not
func (p Password) String() string {
	return string(p)
}

// Encrypt securely encrypts a Password object
func (p Password) Encrypt() Password {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(p.String()), 14)
	return Password(string(bytes))
}

// Check compared password and returns if they match
func (p Password) Check(password Password) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p), []byte(password))
	return err == nil
}
