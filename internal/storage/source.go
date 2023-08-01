package storage

import (
	"context"
	"github.com/HeadGardener/news-feed/internal/models"
	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

type SourceStorage struct {
	db *sqlx.DB
}

func NewSourceStorage(db *sqlx.DB) *SourceStorage {
	return &SourceStorage{db: db}
}

func (s *SourceStorage) Save(ctx context.Context, src models.Source) (int, error) {
	var sourceID int

	err := s.db.QueryRowContext(ctx, saveSourceQuery, src.Name, src.FeedURL, src.CreatedAt).Scan(&sourceID)
	if err != nil {
		return 0, err
	}

	return sourceID, nil
}

func (s *SourceStorage) Sources(ctx context.Context) ([]models.Source, error) {
	var sources []models.Source

	if err := s.db.SelectContext(ctx, &sources, getSourcesQuery); err != nil {
		return nil, err
	}

	return sources, nil
}

func (s *SourceStorage) SourceByID(ctx context.Context, sourceID int) (models.Source, error) {
	var source models.Source

	if err := s.db.GetContext(ctx, &source, getSourceByIDQuery, sourceID); err != nil {
		return models.Source{}, err
	}

	return source, nil
}

func (s *SourceStorage) Delete(ctx context.Context, sourceID int) error {
	if _, err := s.db.ExecContext(ctx, deleteSourceQuery, sourceID); err != nil {
		return err
	}

	return nil
}
