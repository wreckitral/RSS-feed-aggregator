package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    Name string `json:"name"`
}

func NewUser(name string) (*User, error) {
    id := uuid.New()
    
    return &User{
        ID: id,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        Name: name,
    }, nil
}
