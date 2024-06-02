package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/wreckitral/RSS-feed-aggregator/internal/database"
)

type User struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    Name string `json:"name"`
    APIKey string `json:"apiKey"`
}

type CreateUserRequest struct {
    Name string `json:"name"`
}

type Feed struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    Name string `json:"name"`
    Url string `json:"url"`
    UserID uuid.UUID `json:"userId"`
}

type CreateFeedRequest struct {
    Name string `json:"name"`
    Url string `json:"url"`
}

type FeedFollow struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    UserID uuid.UUID `json:"userId"`
    FeedID uuid.UUID `json:"feedId"`
}

type FeedFollowRequest struct {
    FeedID uuid.UUID `json:"feedId"`
}

func NewUser(name string) (*database.User, error) {
    id := uuid.New()
    
    return &database.User{
        ID: id,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name: name,
    }, nil
}

func NewFeed(name, url string, userId uuid.UUID) (*database.Feed, error) {
    id := uuid.New()

    return &database.Feed{
        ID: id,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name: name,
        Url: url,
        UserID: userId,
    }, nil
}

func NewFeedFollow(userId, feedId uuid.UUID) (*database.FeedFollow) {
    id := uuid.New()

    return &database.FeedFollow{
        ID: id,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserID: userId,
        FeedID: feedId,
    }
}
