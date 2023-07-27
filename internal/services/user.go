package services

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/HeadGardener/news-feed/internal/models"
	"time"
)

type UserProcessor interface {
	Create(ctx context.Context, user models.User) (int, error)
	UpdateSendFlag(ctx context.Context, userID, sendFlag int) error
}

type UserService struct {
	userProcessor UserProcessor
}

func NewUserService(userProcessor UserProcessor) *UserService {
	return &UserService{userProcessor: userProcessor}
}

func (s *UserService) Create(ctx context.Context, userInput models.UserInput) (int, error) {
	user := models.User{
		Username:     userInput.Username,
		Email:        userInput.Email,
		PasswordHash: getPasswordHash(userInput.Password),
		Role:         "user",
		SendFlag:     1,
		LastOnline:   time.Now(),
	}

	return s.userProcessor.Create(ctx, user)
}

func (s *UserService) UpdateSendFlag(ctx context.Context, userID, sendFlag int) error {
	return s.UpdateSendFlag(ctx, userID, sendFlag)
}

func getPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
