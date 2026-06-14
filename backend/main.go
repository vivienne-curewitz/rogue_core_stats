package main

import (
	"context"
	"log"
	"sync"

	"github.com/vivienne-curewitz/rogue_core_stats/db"
	"github.com/vivienne-curewitz/rogue_core_stats/handlers"
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
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	log.Printf("Running Rogue Core Stats Server\n")
	dbErr := db.LoadConfig()
	if dbErr != nil {
		log.Fatalf("Failed to load db data: %s\n", dbErr)
	}
	dbErr = db.InitDB(ctx)
	if dbErr != nil {
		log.Fatalf("Failed to init database: %s\n", dbErr)
	}
	wg.Add(1)
	go handlers.StartHandlers(wg)
	wg.Wait()
}
