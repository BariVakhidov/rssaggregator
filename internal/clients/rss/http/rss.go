package rss

import (
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/BariVakhidov/rssaggregator/internal/lib/logger/sl"
	"github.com/BariVakhidov/rssaggregator/internal/model"
)

type Client struct {
	httpClient http.Client
	log        *slog.Logger
}

func New(log *slog.Logger, timeout time.Duration) *Client {
	httpClient := http.Client{
		Timeout: timeout,
	}

	return &Client{
		httpClient: httpClient,
		log:        log,
	}
}

func (c *Client) FetchFeed(feedUrl string) (*model.RSSFeed, error) {
	const op = "client.rss.http.FetchFeed"
	log := c.log.With(slog.String("op", op))

	response, err := c.httpClient.Get(feedUrl)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("failed to close resp body", sl.Err(err))
		}
	}(response.Body)

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	rssFeed := model.RSSFeed{}

	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return nil, err
	}

	return &rssFeed, nil
}
