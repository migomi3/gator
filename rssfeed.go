package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/migomi3/gator/internal/database"
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
		parsedDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return nil, err
		}
		pubDate := sql.NullTime{
			Time:  parsedDate,
			Valid: true,
		}

		params := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
			},
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		}

		post, err := s.db.CreatePost(context.Background(), params)
		if err != nil {
			return nil, err
		}

		fmt.Println(post)
	}

	return rssFeed, nil
}
