package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type FavoritesStorage struct {
	db *sqlx.DB
}

func NewFavoritesStorage(db *sqlx.DB) *FavoritesStorage {
	return &FavoritesStorage{db: db}
}

func (s *FavoritesStorage) Add(ctx context.Context, userID, articleID int) error {
	if _, err := s.db.ExecContext(ctx, addFavoriteQuery, userID, articleID); err != nil {
		return err
	}

	return nil
}

func (s *FavoritesStorage) ArticlesByUserID(ctx context.Context, userID int) ([]int, error) {
	var articlesIDs []int

	if err := s.db.SelectContext(ctx, &articlesIDs, getArticlesByUserIDQuery, userID); err != nil {
		return nil, err
	}

	return articlesIDs, nil
}

func (s *FavoritesStorage) Delete(ctx context.Context, userID, articleID int) error {
	if _, err := s.db.ExecContext(ctx, deleteFavoriteQuery, userID, articleID); err != nil {
		return err
	}

	return nil
}
