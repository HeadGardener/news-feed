package fetcher

import (
	"context"
	"github.com/HeadGardener/news-feed/internal/lib/source"
	"github.com/HeadGardener/news-feed/internal/models"
	"log"
	"strings"
	"sync"
	"time"
)

const fetchInterval = time.Minute

var categories = []string{"go", "golang"}

type ArticleSaver interface {
	Save(ctx context.Context, art models.Article) error
}

type SourceProvider interface {
	Sources(ctx context.Context) ([]models.Source, error)
}

type Source interface {
	Fetch(ctc context.Context) ([]models.Item, error)
	ID() int
	Name() string
}

type Fetcher struct {
	sourceProvider SourceProvider
	articleSaver   ArticleSaver
}

func NewFetcher(sourceProvider SourceProvider, articleSaver ArticleSaver) *Fetcher {
	return &Fetcher{
		sourceProvider: sourceProvider,
		articleSaver:   articleSaver,
	}
}

func (f *Fetcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(fetchInterval)
	defer ticker.Stop()

	if err := f.Fetch(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := f.Fetch(ctx); err != nil {
				return err
			}
			log.Printf("[]")
		}
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	sources, err := f.sourceProvider.Sources(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, src := range sources {
		wg.Add(1)

		go func(source Source) {
			items, err := source.Fetch(ctx)
			if err != nil {
				log.Printf("[ERROR] failed to fetch items: %e", err)
				return
			}

			if err := f.handleItems(ctx, items, source); err != nil {
				log.Printf("[ERROR] failed while handling items items: %e", err)
				return
			}
			wg.Done()
		}(source.NewRSSSource(src))
	}
	wg.Wait()

	return nil
}

func (f *Fetcher) handleItems(ctx context.Context, items []models.Item, source Source) error {
	for _, item := range items {
		if ignoreItem(item) {
			log.Printf("[INFO] item %s from %s was ignored", item.Link, source.Name())
			continue
		}

		art := models.Article{
			SourceID:    source.ID(),
			Title:       item.Title,
			Link:        item.Link,
			Summary:     item.Summary,
			PublishedAt: item.Date.UTC(),
		}

		if err := f.articleSaver.Save(ctx, art); err != nil {
			return err
		}
	}

	return nil
}

func ignoreItem(item models.Item) bool {
	for _, c := range item.Categories {
		cl := strings.ToLower(c)
		if strings.Compare(cl, categories[0]) == 0 || strings.Compare(cl, categories[1]) == 0 {
			return false
		}
	}
	return true
}
