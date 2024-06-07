package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wreckitral/RSS-feed-aggregator/internal/database"
)

func startScrapping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
    log.Printf("Scrapping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
    ticker := time.NewTicker(timeBetweenRequest)
    for ; ; <-ticker.C {
        feeds, err := db.GetNextFeedsToFetch(
            context.Background(), 
            int32(concurrency),
        )
        if err != nil {
            log.Println("error fetching feeds:", err)
            continue
        }

        wg := &sync.WaitGroup{}
        for _, feed := range feeds {
            wg.Add(1)

            go scrapeFeed(db, wg, feed)
        }
        wg.Wait()
    }
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
    defer wg.Done()

    _, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
    if err != nil {
        log.Println("Error marking feed as fetched:", err)
        return 
    }

    rssFeed, err := urlToFeed(feed.Url)
    if err != nil {
        log.Println("Error fetching feed:", err)
    }

    for _, item := range rssFeed.Channel.Item {
        desc := sql.NullString{}
        if item.Description != "" {
            desc.String = item.Description
            desc.Valid = true
        }

        pubAt, err := parseTime(item.PubDate)
        if err != nil {
            log.Println("error parsing time:", err, "on item", item.Title)
        }

        _, err = db.CreatePost(context.Background(), 
            database.CreatePostParams{
                ID: uuid.New(),
                CreatedAt: time.Now().UTC(),
                UpdatedAt: time.Now().UTC(),
                Title: item.Title,
                Url: item.Link,
                Description: desc,
                PublishedAt: pubAt,
                FeedID: feed.ID,
            },
        )
        if err != nil {
            if strings.Contains(err.Error(), "duplicate key") {
                continue
            }
            log.Println("failed to create new post:", err)
        }
    }

    log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

func parseTime(timeStr string) (time.Time, error) {
    // Primary format to try
    formats := []string{
        time.RFC1123Z,
        time.RFC1123,
        time.RFC3339,
    }

    var t time.Time
    var err error

    for _, format := range formats {
        t, err = time.Parse(format, timeStr)
        if err == nil {
            return t, nil
        }
    }
    return time.Time{}, fmt.Errorf("unable to parse time: %v", timeStr)
}
