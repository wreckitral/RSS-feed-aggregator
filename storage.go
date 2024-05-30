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
    CreateUserToDb(*User) error
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

func (s *PostgresStore) CreateUserToDb(user *User) error {
    userToDb := database.CreateUserParams{
        ID: user.ID,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
        Name: user.Name,
    }

    _, err := s.DB.CreateUser(context.Background(), userToDb) 
    if err != nil {
        return err
    }

    return nil
}
