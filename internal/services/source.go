package services

import (
	"context"
	"github.com/HeadGardener/news-feed/internal/models"
)

type SourceProcessor interface {
	Save(ctx context.Context, src models.Source) (int, error)
	Delete(ctx context.Context, sourceID int) error
}

type SourceService struct {
	sourceProcessor SourceProcessor
}

func NewSourceService(sourceSaver SourceProcessor) *SourceService {
	return &SourceService{sourceProcessor: sourceSaver}
}

func (s *SourceService) Save(ctx context.Context, srcInput models.SourceInput) (int, error) {
	src := models.Source{
		Name:      srcInput.Name,
		FeedURL:   srcInput.FeedURL,
		CreatedAt: srcInput.CreatedAt,
	}

	return s.sourceProcessor.Save(ctx, src)
}

func (s *SourceService) Delete(ctx context.Context, sourceID int) error {
	return s.sourceProcessor.Delete(ctx, sourceID)
}
