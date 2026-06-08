package main

import (
	"context"
	"log"

	"github.com/vivienne-curewitz/rogue_core_stats/db"
)

func initDb() {
	ctx := context.Background()
	log.Printf("Initializing database connection\n")
	if err := db.LoadConfig(); err != nil {
		log.Fatalf("Failed to load db data: %s\n", err)
	}
	if err := db.InitDB(ctx); err != nil {
		log.Fatalf("Failed to init database: %s\n", err)
	}
}

func main() {
	log.Printf("Running Rogue Core Stats Server\n")
	initDb()
}
