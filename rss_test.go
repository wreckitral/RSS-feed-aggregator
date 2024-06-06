package main

import (
	"testing"
)

func TestRss(t *testing.T) {
	feed, err := urlToFeed("https://wagslane.dev/index.xml")
	if err != nil {
		t.Fatalf("Failed to get feed: %v", err)
	}

	// Add relevant checks for the fields in RSSFeed struct
	expectedTitle := "Lane's Blog"
	expectedDescription := "Recent content on Lane's Blog"

	if feed.Channel.Title != expectedTitle {
		t.Errorf("Expected title %s, but got %s", expectedTitle, feed.Channel.Title)
	}

	if feed.Channel.Description != expectedDescription {
		t.Errorf("Expected description %s, but got %s", expectedDescription, feed.Channel.Description)
	}

	t.Logf("Feed: %+v", feed)
}

