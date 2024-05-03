package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(feedUrl string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := httpClient.Get(feedUrl)

	if err != nil {
		return RSSFeed{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(response.Body)

	data, readErr := io.ReadAll(response.Body)

	if readErr != nil {
		return RSSFeed{}, readErr
	}

	rssFeed := RSSFeed{}

	xmlErr := xml.Unmarshal(data, &rssFeed)
	if xmlErr != nil {
		return RSSFeed{}, xmlErr
	}

	return rssFeed, nil
}
