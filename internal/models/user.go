package models

import (
	"errors"
	"regexp"
	"time"
)

var (
	checkPassword = regexp.MustCompile(`[0-9A-z]{8,16}$`)
	checkEmail    = regexp.MustCompile(`[0-9A-z]@gmail.com$`)
)

type User struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	SendFlag     int       `db:"send_flag"`
	LastOnline   time.Time `db:"last_online"`
}

type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserAttributes struct {
	ID    int    `json:"id"`
	Email string `json:"name"`
	Role  string `json:"role"`
}

func (u *UserInput) Validate() error {
	if !checkEmail.MatchString(u.Email) {
		return errors.New("invalid email (only gmail available)")
	}
	if !checkPassword.MatchString(u.Password) {
		return errors.New("invalid password (length 8-16, only letters and numbers)")
	}

	return nil
}
