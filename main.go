package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(".env"); err != nil {
        log.Fatal(err)
    }

    port := ":" + os.Getenv("PORT")

    store, err := NewPostgresStore()
    if err != nil {
        log.Fatal(err)
    }

    server := NewAPIServer(port, store)

    if err := server.Run(); err != nil {
        log.Fatal(err)
    }
}
