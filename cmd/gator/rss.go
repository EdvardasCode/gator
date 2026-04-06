package main

import (
	"context"
	"encoding/xml"
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

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	request.Header.Set("User-Agent", "gator")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return &RSSFeed{}, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	var rssFeed RSSFeed
	if err = xml.Unmarshal(data, &rssFeed); err != nil {
		return &RSSFeed{}, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
	}

	return &rssFeed, nil
}
