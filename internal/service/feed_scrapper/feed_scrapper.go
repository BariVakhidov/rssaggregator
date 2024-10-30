package feedscrapper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/BariVakhidov/rssaggregator/internal/lib/logger/sl"
	"github.com/BariVakhidov/rssaggregator/internal/model"
	"github.com/BariVakhidov/rssaggregator/internal/storage"
	"github.com/google/uuid"
)

type FeedProvider interface {
	NextFeedsToFetch(ctx context.Context, limit int) ([]model.Feed, error)
	MarkFeedAsFetched(ctx context.Context, feedId uuid.UUID) (*model.Feed, error)
}

type PostProvider interface {
	CreatePost(ctx context.Context, postInfo model.CreatePostInfo) (*model.Post, error)
}

type FeedFetcher interface {
	FetchFeed(feedUrl string) (*model.RSSFeed, error)
}

type Scrapper struct {
	log          *slog.Logger
	postProvider PostProvider
	feedProvider FeedProvider
	feedFetcher  FeedFetcher
}

func New(log *slog.Logger, postProvider PostProvider, feedProvider FeedProvider, feedFetcher FeedFetcher) *Scrapper {
	return &Scrapper{
		log:          log,
		feedProvider: feedProvider,
		postProvider: postProvider,
		feedFetcher:  feedFetcher,
	}
}

func (s *Scrapper) RunScrapper(concurrency int, timeBetweenRequest time.Duration) {
	const op = "service.feed_scrapper.StartScrapping"
	log := s.log.With(slog.String("op", op))

	log.Info(fmt.Sprintf("Scrapping on %d goroutines every %s duration", concurrency, timeBetweenRequest))

	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := s.feedProvider.NextFeedsToFetch(context.Background(), concurrency)

		if err != nil {
			log.Error("error fetching feeds", sl.Err(err))
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go s.scrapeFeed(wg, feed)
		}

		wg.Wait()
	}
}

func (s *Scrapper) scrapeFeed(wg *sync.WaitGroup, feed model.Feed) {
	const op = "service.feed_scrapper.StartScrapping"
	log := s.log.With(slog.String("op", op))

	defer wg.Done()

	_, err := s.feedProvider.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Error("error marking feed as fetched", sl.Err(err))
		return
	}

	rssFeeds, err := s.feedFetcher.FetchFeed(feed.Url)
	if err != nil {
		log.Error("error fetching feed", slog.String("url", feed.Url), sl.Err(err))
		return
	}

	for _, rssFeed := range rssFeeds.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123, rssFeed.PubDate)
		if err != nil {
			log.Error("couldn't parse date", sl.Err(err))
			continue
		}

		_, err = s.postProvider.CreatePost(context.Background(), model.CreatePostInfo{
			ID:          uuid.New(),
			FeedID:      feed.ID,
			Title:       rssFeed.Title,
			Url:         rssFeed.Link,
			Description: rssFeed.Description,
			PublishedAt: publishedAt,
		})

		if err != nil {
			if errors.Is(err, storage.ErrPostExists) {
				continue
			}
			log.Error("creating post failed", slog.String("rssFeed.Title", rssFeed.Title), sl.Err(err))
		}
	}

	log.Info(fmt.Sprintf("Feed %s collected, %v posts found", feed.Name, len(rssFeeds.Channel.Item)))
}
