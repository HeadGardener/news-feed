package storage

import (
	"context"
	"github.com/HeadGardener/news-feed/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) Create(ctx context.Context, user models.User) (int, error) {
	var userID int

	err := s.db.QueryRowContext(ctx, createUserQuery, user.Username, user.Email, user.PasswordHash).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *UserStorage) Users(ctx context.Context) ([]models.User, error) {
	var users []models.User

	if err := s.db.SelectContext(ctx, &users, getUsersForSendQuery); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserStorage) UserByInput(ctx context.Context, userInput models.UserInput) (models.User, error) {
	var user models.User

	err := s.db.GetContext(ctx, &user, getUserWithInputQuery, userInput.Username, userInput.Email, userInput.Password)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *UserStorage) UpdateSendFlag(ctx context.Context, userID, sendFlag int) error {
	if _, err := s.db.ExecContext(ctx, updateSendFlagQuery, sendFlag, userID); err != nil {
		return err
	}

	return nil
}
