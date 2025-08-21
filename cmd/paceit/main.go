package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lorenas/PaceIt/internal/db"
)

func main() {
    _ = godotenv.Load()
    conn, err := db.Open()
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Println("db connection ok, placeholder app on port", port)
}
