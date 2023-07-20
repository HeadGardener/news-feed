package services

import (
	"context"
	"github.com/HeadGardener/news-feed/internal/models"
)

type ArticleProvider interface {
	Articles(ctx context.Context, page int) ([]models.Article, error)
}

type ArticleService struct {
	articleProvider ArticleProvider
}

func NewArticleService(articleProvider ArticleProvider) *ArticleService {
	return &ArticleService{
		articleProvider: articleProvider,
	}
}

func (s *ArticleService) GetAll(ctx context.Context, page int) ([]models.Article, error) {
	return s.articleProvider.Articles(ctx, page)
}
