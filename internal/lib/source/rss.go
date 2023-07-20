package source

import (
	"context"
	"fmt"
	"github.com/HeadGardener/news-feed/internal/models"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
)

type RSSSource struct {
	URL        string
	SourceID   int
	SourceName string
}

func NewRSSSource(s models.Source) RSSSource {
	return RSSSource{
		URL:        s.FeedURL,
		SourceID:   s.ID,
		SourceName: s.Name,
	}
}

func (s RSSSource) Fetch(ctx context.Context) ([]models.Item, error) {
	fmt.Println(s.URL)
	feed, err := s.loadFeed(ctx, s.URL)
	if err != nil {
		return nil, err
	}

	return lo.Map(feed.Items, func(item *gofeed.Item, _ int) models.Item {
		return models.Item{
			Title:      item.Title,
			Categories: item.Categories,
			Link:       item.Link,
			Date:       *item.PublishedParsed,
			Summary:    item.Description,
			SourceName: s.SourceName,
		}
	}), nil
}

func (s RSSSource) loadFeed(ctx context.Context, url string) (*gofeed.Feed, error) {
	var (
		feedCh = make(chan *gofeed.Feed)
		errCh  = make(chan error)
		fp     = gofeed.NewParser()
	)

	go func() {
		feed, err := fp.ParseURL(url)
		if err != nil {
			errCh <- err
			return
		}

		feedCh <- feed
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case feed := <-feedCh:
		return feed, nil
	}
}

func (s RSSSource) ID() int {
	return s.SourceID
}

func (s RSSSource) Name() string {
	return s.SourceName
}
