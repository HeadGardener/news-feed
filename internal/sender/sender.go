package sender

import (
	"context"
	"github.com/HeadGardener/news-feed/internal/lib/email"
	"github.com/HeadGardener/news-feed/internal/models"
	"log"
	"sync"
	"time"
)

const sendInterval = 24 * time.Hour

type UserProvider interface {
	Users(ctx context.Context) ([]models.User, error)
}

type Sender struct {
	userProvider UserProvider
}

func NewSender(userProvider UserProvider) *Sender {
	return &Sender{userProvider: userProvider}
}

func (s *Sender) Start(ctx context.Context) error {
	ticker := time.NewTicker(sendInterval)
	defer ticker.Stop()

	if err := s.Process(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.Process(ctx); err != nil {
				return err
			}
			log.Printf("[]")
		}
	}
}

func (s *Sender) Process(ctx context.Context) error {
	email.InitEmail()

	users, err := s.userProvider.Users(ctx)
	if err != nil {
		log.Printf("[ERROR] unable to get users: %e", err)
		return err
	}

	var wg sync.WaitGroup

	for _, user := range users {
		wg.Add(1)

		go func(u models.User) {
			ulo := u.LastOnline.Add(72 * time.Hour)
			if ulo.Before(time.Now().Add(3 * time.Hour).UTC()) {
				return
			}

			if err := email.Send(u.Email); err != nil {
				log.Printf("[ERROR] failed to send email for user %s (%d): %e", u.Username, u.ID, err)
			}

			wg.Done()
		}(user)
	}

	wg.Wait()

	return nil
}
