package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	if err = xml.Unmarshal(body, &rssFeed); err != nil {
		return nil, err
	}

	return unEscapeFeed(&rssFeed), nil
}

func unEscapeFeed(feed *RSSFeed) *RSSFeed {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		item.Title = html.EscapeString(item.Title)
		item.Description = html.EscapeString(item.Description)
	}

	return feed
}

func scrapeFeeds(s *state) (*RSSFeed, error) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return nil, err
	}

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return nil, err
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return nil, err
	}

	for _, item := range rssFeed.Channel.Item {
		fmt.Println(item.Title)
	}

	return rssFeed, nil
}
