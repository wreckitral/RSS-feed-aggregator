package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/wreckitral/RSS-feed-aggregator/internal/database"
)

// Storage is the interface that stores logic for interacting with database
// and marshalling data from database to json for API response
type Storage interface {
    CreateUserToDb(*database.User)                  (*User, error)
    GetUserByAPIKey(apiKey string)                  (*User, error)
    CreateFeedToDb(*database.Feed)                  (*Feed, error)
    GetFeeds()                                      ([]*Feed, error)  
    CreateFeedFollows(*database.FeedFollow)         (*FeedFollow, error)
    GetFeedFollows(uuid.UUID)                       ([]*FeedFollow, error)
    DeleteFeedFollows(id, user_id uuid.UUID)        error
    GetPostForUsers(userId uuid.UUID, limit int32)  ([]*Post, error)
}

type PostgresStore struct {
    DB *database.Queries
}

func NewPostgresStore() (*PostgresStore, error) {
    dbConn := os.Getenv("DBCONN")
    if dbConn == "" {
        log.Fatal("db env is missing")
    }

    db, err := sql.Open("postgres", dbConn)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    dbQueries := database.New(db)
    
    go startScrapping(dbQueries, 10, time.Minute)
    
    return &PostgresStore{
        DB: dbQueries,
    }, nil
}

func (s *PostgresStore) CreateUserToDb(user *database.User) (*User, error){
    if user.Name == "" {
        return nil, fmt.Errorf("name cannot be empty string")
    }

    if len(user.Name) <= 3{
        return nil, fmt.Errorf("name is too short")
    }
    
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
    if feed.Name == "" || feed.Url == "" {
        return nil, fmt.Errorf("request data cannot be empty")
    }

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

func (s *PostgresStore) CreateFeedFollows(ff *database.FeedFollow) (*FeedFollow, error) {
    feedToDb := database.CreateFeedFollowsParams{
        ID: ff.ID,
        CreatedAt: ff.CreatedAt,
        UpdatedAt: ff.UpdatedAt,
        UserID: ff.UserID,
        FeedID: ff.FeedID,
    }

    feedToAPI, err := s.DB.CreateFeedFollows(context.Background(), feedToDb)
    if err != nil {
        return nil, err
    }

    return &FeedFollow{
        ID: feedToAPI.ID,
        CreatedAt: feedToAPI.CreatedAt,
        UpdatedAt: feedToAPI.UpdatedAt,
        UserID: feedToAPI.UserID,
        FeedID: feedToAPI.FeedID,
    }, nil
}

func (s *PostgresStore) DeleteFeedFollows(id, userId uuid.UUID) error {
    feedFollowsToDelete := database.DeleteFeedFollowsParams{
        ID: id,
        UserID: userId,
    }

    if err := s.DB.DeleteFeedFollows(context.Background(), feedFollowsToDelete); err != nil {
        return err
    }

    return nil
}

func (s *PostgresStore) GetFeedFollows(userId uuid.UUID) ([]*FeedFollow, error) {
    getFeedFollows, err := s.DB.GetAllFeedFollows(context.Background(), userId)
    if err != nil {
        return nil, err
    }

    feedFollows := []*FeedFollow{}

    for _, feed := range getFeedFollows {
        f := &FeedFollow {
            ID: feed.ID,
            FeedID: feed.FeedID,
            UserID: feed.UserID,
            CreatedAt: feed.CreatedAt,
            UpdatedAt: feed.UpdatedAt,
        }
        feedFollows = append(feedFollows, f)
    }

    return feedFollows, nil
}


func (s *PostgresStore) GetPostForUsers(userId uuid.UUID, limit int32) ([]*Post, error) {
    params := database.GetPostsForUserParams{
        UserID: userId,        
        Limit: limit,
    }

    getPosts, err := s.DB.GetPostsForUser(context.Background(), params)
    if err != nil {
        return nil, err
    }

    posts := []*Post{}

    for _, post := range getPosts {
        var description *string
        if post.Description.Valid {
            description = &post.Description.String
        }
        p := &Post{
            ID: post.ID,
            CreatedAt: post.CreatedAt,
            UpdatedAt: post.UpdatedAt,
            Title: post.Title,
            Url: post.Url,
            Description: description,
            PublishedAt: post.PublishedAt,
            FeedID: post.FeedID,
        }
        posts = append(posts, p)
    }

    return posts, nil
}
