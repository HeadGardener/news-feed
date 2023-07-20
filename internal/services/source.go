package services

import (
	"context"
	"github.com/HeadGardener/news-feed/internal/models"
)

type SourceSaver interface {
	Save(ctx context.Context, src models.Source) (int, error)
}

type SourceService struct {
	sourceSaver SourceSaver
}

func NewSourceService(sourceSaver SourceSaver) *SourceService {
	return &SourceService{sourceSaver: sourceSaver}
}

func (s *SourceService) Save(ctx context.Context, srcInput models.SourceInput) (int, error) {
	src := models.Source{
		Name:      srcInput.Name,
		FeedURL:   srcInput.FeedURL,
		CreatedAt: srcInput.CreatedAt,
	}

	return s.sourceSaver.Save(ctx, src)
}
