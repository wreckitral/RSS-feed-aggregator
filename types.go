package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/wreckitral/RSS-feed-aggregator/internal/database"
)

type UserResponse struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    Name string `json:"name"`
    APIKey string `json:"apiKey"`
}

type CreateUserRequest struct {
    Name string `json:"name"`
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
