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

    server := NewAPIServer(port)

    if err := server.Run(); err != nil {
        log.Fatal(err)
    }
}
