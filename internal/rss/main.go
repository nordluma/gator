package rss

import (
	"context"
	"encoding/xml"
	"html"
	"net/http"
	"time"
)

func FetchFeed(ctx context.Context, feedUrl string) (*Feed, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return &Feed{}, err
	}
	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return &Feed{}, err
	}
	defer res.Body.Close()

	var feed Feed
	decoder := xml.NewDecoder(res.Body)
	if err = decoder.Decode(&feed); err != nil {
		return &Feed{}, err
	}

	escapeHtml(&feed)

	return &feed, nil
}

func escapeHtml(feed *Feed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}
}
