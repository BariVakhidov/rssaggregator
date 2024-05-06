package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/BariVakhidov/rssaggregator/internal/database"
	"github.com/google/uuid"
)

func startScrapping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scrapping on %d goroutines every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))

		if err != nil {
			log.Println("error fetching feeds: ", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(wg, feed, db)
		}

		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup, feed database.Feed, db *database.Queries) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("error marking feed as fetched ", err)
		return
	}

	rssFeeds, rssErr := fetchFeed(feed.Url)
	if rssErr != nil {
		log.Printf("error fetching %s : %v", feed.Url, rssErr)
		return
	}

	for _, rssFeed := range rssFeeds.Channel.Item {
		description := sql.NullString{}
		if rssFeed.Description != "" {
			description.String = rssFeed.Description
			description.Valid = true
		}

		publishedAt, timeErr := time.Parse(time.RFC1123, rssFeed.PubDate)

		if timeErr != nil {
			log.Printf("couldn't parse date %v with err %v", rssFeed.PubDate, timeErr)
			continue
		}

		_, postErr := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       rssFeed.Title,
			Description: description,
			FeedID:      feed.ID,
			PublishedAt: publishedAt,
			Url:         rssFeed.Link,
		})

		if postErr != nil {
			if strings.Contains(postErr.Error(), "duplicate key") {
				continue
			}
			log.Printf("error creating post %s : %v", rssFeed.Title, postErr)
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeeds.Channel.Item))
}
