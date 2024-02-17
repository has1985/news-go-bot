package source

import (
	"context"
	"fmt"
	"github.com/SlyMarbo/rss"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/has1985/news-go-bot/internal/model"
)

type RSSSource struct {
	URL        string
	SourceID   int64
	SourceName string
}

func NewRSSSourceFromModel(m model.Source) *RSSSource {
	return &RSSSource{
		URL:        m.FeedURL,
		SourceID:   m.ID,
		SourceName: m.Name,
	}
}

func (s *RSSSource) Fetch(ctx context.Context) ([]model.Item, error) {
	logger := ctxlogrus.Extract(ctx)
	logger.Debugln("RSSSource.Fetch")

	feed, err := s.loadFeed(ctx, s.URL)
	if err != nil {
		logger.WithError(err).Errorf("Fetch: unable to load feed. SourceID: %d", s.SourceID)
		return nil, fmt.Errorf("unable to load feed")
	}

	var items []model.Item
	for _, feedItem := range feed.Items {
		item := model.Item{
			Title:      feedItem.Title,
			Categories: feedItem.Categories,
			Link:       feedItem.Link,
			Data:       feedItem.Date,
			Summary:    feedItem.Summary,
			SourceName: s.SourceName,
		}

		items = append(items, item)
	}

	return items, nil
}

func (s *RSSSource) loadFeed(ctx context.Context, url string) (*rss.Feed, error) {
	logger := ctxlogrus.Extract(ctx)
	logger.Debugln("loadFeed")

	var (
		feedCh = make(chan *rss.Feed)
		errCh  = make(chan error)
	)

	go func() {
		feed, err := rss.Fetch(url)
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
	case feedCh := <-feedCh:
		return feedCh, nil
	}

}

func (s *RSSSource) ID() int64 {
	return s.SourceID
}

func (s *RSSSource) Name() string {
	return s.SourceName
}
