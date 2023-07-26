package services

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/HeadGardener/news-feed/internal/models"
	"time"
)

type UserProvider interface {
	Create(ctx context.Context, user models.User) (int, error)
}

type UserService struct {
	userProcessor UserProvider
}

func NewUserService(userProcessor UserProvider) *UserService {
	return &UserService{userProcessor: userProcessor}
}

func (s *UserService) Create(ctx context.Context, userInput models.UserInput) (int, error) {
	user := models.User{
		Username:     userInput.Username,
		Email:        userInput.Email,
		PasswordHash: getPasswordHash(userInput.Password),
		SendFlag:     1,
		LastOnline:   time.Now(),
	}

	return s.userProcessor.Create(ctx, user)
}

func getPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
