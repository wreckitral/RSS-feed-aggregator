package main

import (
	"log"
	"os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        log.Fatal("port env is missing")
    }

    store, err := NewPostgresStore()
    if err != nil {
        log.Fatal(err)
    }

    server := NewAPIServer(port, store)

    if err := server.Run(); err != nil {
        log.Fatal(err)
    }
}
