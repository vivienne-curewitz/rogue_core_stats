package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig loads the database configuration from .env.config and
// sets the DATABASE_URL environment variable used by InitDB.
func LoadConfig() error {
	// Load .env.config file if it exists
	if err := godotenv.Load(".env.config"); err != nil {
		// It's okay if the file doesn't exist, as long as the variables are set in the environment
		fmt.Printf("Warning: .env.config not found or could not be loaded: %v\n", err)
	}

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	if user == "" || pass == "" || host == "" || port == "" || name == "" {
		return fmt.Errorf("missing required database environment variables (DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)")
	}

	// Construct connection string: postgres://user:pass@host:port/dbname
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, name)
	
	// Set DATABASE_URL for InitDB to consume
	os.Setenv("DATABASE_URL", connStr)

	return nil
}
