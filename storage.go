package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/wreckitral/RSS-feed-aggregator/internal/database"
)

type Storage interface {
    CreateUserToDb(*database.User)  (*User, error)
    GetUserByAPIKey(apiKey string)  (*User, error)
    CreateFeedToDb(*database.Feed)  (*Feed, error)
    GetFeeds()                      ([]*Feed, error)  
}

type PostgresStore struct {
    DB *database.Queries
}

func NewPostgresStore() (*PostgresStore, error) {
    dbConn := os.Getenv("DBCONN")
    fmt.Println(dbConn)

    db, err := sql.Open("postgres", dbConn)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    dbQueries := database.New(db)
    
    return &PostgresStore{
        DB: dbQueries,
    }, nil
}

func (s *PostgresStore) CreateUserToDb(user *database.User) (*User, error){
    userToDb := database.CreateUserParams{
        ID: user.ID,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
        Name: user.Name,
    }

    createdUser, err := s.DB.CreateUser(context.Background(), userToDb) 
    if err != nil {
        return nil, err
    }

    userToAPI := &User{
        ID: createdUser.ID,
        CreatedAt: createdUser.CreatedAt,
        UpdatedAt: createdUser.UpdatedAt,
        Name: createdUser.Name,
        APIKey: createdUser.ApiKey,
    }

    return userToAPI, nil
}

func (s *PostgresStore) GetUserByAPIKey(apiKey string) (*User, error) {
    user, err := s.DB.GetUserByAPIKey(context.Background(), apiKey)
    if err != nil {
        return nil, err
    }

    userToAPI := &User{
        ID: user.ID,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
        Name: user.Name,
        APIKey: user.ApiKey,
    }
    
    return userToAPI, nil
}


func (s *PostgresStore) CreateFeedToDb(feed *database.Feed) (*Feed, error) {
    feedToDb := database.CreateFeedParams{
        ID: feed.ID,
        CreatedAt: feed.CreatedAt,
        UpdatedAt: feed.UpdatedAt,
        Name: feed.Name,
        Url: feed.Url,
        UserID: feed.UserID,
    }

    createdFeed, err := s.DB.CreateFeed(context.Background(), feedToDb)
    if err != nil {
        return nil, err
    }

    feedToAPI := Feed{
        ID: createdFeed.ID,
        CreatedAt: createdFeed.CreatedAt,
        UpdatedAt: createdFeed.UpdatedAt,
        Name: createdFeed.Name,
        Url: createdFeed.Url,
        UserID: createdFeed.UserID,
    }

    return &feedToAPI, nil
}

func (s *PostgresStore) GetFeeds() ([]*Feed, error) {
    feedsFromDb, err := s.DB.GetAllFeeds(context.Background())
    if err != nil {
        return nil, err
    }

    feeds := []*Feed{}

    for _, feed := range feedsFromDb {
        f := &Feed {
            ID: feed.ID,
            CreatedAt: feed.CreatedAt,
            UpdatedAt: feed.UpdatedAt,
            Name: feed.Name,
            Url: feed.Url,
            UserID: feed.UserID,
        }
        feeds = append(feeds, f)
    }

    return feeds, nil
}
