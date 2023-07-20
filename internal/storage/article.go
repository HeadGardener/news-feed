package storage

import (
	"context"
	"github.com/HeadGardener/news-feed/internal/models"
	"github.com/jmoiron/sqlx"
)

type ArticleStorage struct {
	db *sqlx.DB
}

func NewArticleStorage(db *sqlx.DB) *ArticleStorage {
	return &ArticleStorage{db: db}
}

func (s *ArticleStorage) Save(ctx context.Context, art models.Article) error {
	_, err := s.db.ExecContext(ctx,
		saveArticleQuery,
		art.SourceID,
		art.Title,
		art.Link,
		art.Summary,
		art.PublishedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *ArticleStorage) Articles(ctx context.Context, page int) ([]models.Article, error) {
	var articles []models.Article

	if err := s.db.SelectContext(ctx, &articles, getArticlesQuery((page-1)*10+1, page*10)); err != nil {
		return nil, err
	}

	return articles, nil
}
