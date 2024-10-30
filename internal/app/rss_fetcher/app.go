package rssfetcherapp

import (
	"log/slog"
	"time"

	rss "github.com/BariVakhidov/rssaggregator/internal/clients/rss/http"
	feedscrapper "github.com/BariVakhidov/rssaggregator/internal/service/feed_scrapper"
)

type App struct {
	scrapper *feedscrapper.Scrapper
}

func New(log *slog.Logger,
	postProvider feedscrapper.PostProvider,
	feedProvider feedscrapper.FeedProvider,
	timeout time.Duration,
) *App {
	feedFetcher := rss.New(log, timeout)
	feedScrapper := feedscrapper.New(log, postProvider, feedProvider, feedFetcher)

	return &App{scrapper: feedScrapper}
}

func (a *App) RunScrapper(concurrency int, timeBetweenRequest time.Duration) {
	a.scrapper.RunScrapper(concurrency, timeBetweenRequest)
}
