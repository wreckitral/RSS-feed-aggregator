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
    CreateUserToDb(*database.User) (*database.User, error)
    GetUserByAPIKey(apiKey string) (*database.User, error)
    CreateFeedToDb(*database.Feed) (*database.Feed, error)
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

func (s *PostgresStore) CreateUserToDb(user *database.User) (*database.User, error){
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

    return &createdUser, nil
}

func (s *PostgresStore) GetUserByAPIKey(apiKey string) (*database.User, error) {
    user, err := s.DB.GetUserByAPIKey(context.Background(), apiKey)
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}


func (s *PostgresStore) CreateFeedToDb(feed *database.Feed) (*database.Feed, error) {
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

    return &createdFeed, nil
}


