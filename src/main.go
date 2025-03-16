package main

import (
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

func main() {
	// Only load .env if not running in Docker (i.e., env vars are not set)
	if os.Getenv("DB_PATH") == "" || os.Getenv("PORT") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("No DB_PATH environment variable")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("No PORT environment variable")
	}

	store, err := NewSQLiteStorage(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	server := NewServer(port, store)
	server.Run()
}
