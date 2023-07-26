package services

import (
	"context"
	"github.com/HeadGardener/news-feed/internal/models"
)

type FavoritesProcessor interface {
	Add(ctx context.Context, userID, articleID int) error
	ArticlesByUserID(ctx context.Context, userID int) ([]int, error)
	Delete(ctx context.Context, userID, articleID int) error
}

type ArticleGetter interface {
	ArticleByID(ctx context.Context, articleID int) (models.Article, error)
}

type FavoritesService struct {
	favoritesProcessor FavoritesProcessor
	articleGetter      ArticleGetter
}

func NewFavoritesService(favoritesProcessor FavoritesProcessor, articleGetter ArticleGetter) *FavoritesService {
	return &FavoritesService{
		favoritesProcessor: favoritesProcessor,
		articleGetter:      articleGetter,
	}
}

func (s *FavoritesService) Add(ctx context.Context, userID, articleID int) error {
	if _, err := s.articleGetter.ArticleByID(ctx, articleID); err != nil {
		return err
	}

	return s.favoritesProcessor.Add(ctx, userID, articleID)
}

func (s *FavoritesService) GetAll(ctx context.Context, userID int) ([]models.Article, error) {
	articlesIDs, err := s.favoritesProcessor.ArticlesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var articles []models.Article
	for _, id := range articlesIDs {
		a, err := s.articleGetter.ArticleByID(ctx, id)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, nil
}

func (s *FavoritesService) Delete(ctx context.Context, userID, articleID int) error {
	return s.favoritesProcessor.Delete(ctx, userID, articleID)
}
