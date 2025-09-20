package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/manuel/make-it-rain/config"
	"github.com/manuel/make-it-rain/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if err := config.LoadConfig("."); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run migrate.go [up|down] [steps]")
	}

	command := os.Args[1]
	databaseURL := config.Cfg.Database.GetConnectionString()

	// Debug output
	fmt.Printf("Database connection: host=%s, port=%d, user=%s, dbname=%s\n",
		config.Cfg.Database.Host,
		config.Cfg.Database.Port,
		config.Cfg.Database.User,
		config.Cfg.Database.Name)

	switch command {
	case "up":
		if err := db.RunMigrations(databaseURL); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully")

	case "down":
		steps := 1
		if len(os.Args) > 2 {
			var err error
			steps, err = strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalf("Invalid number of steps: %v", err)
			}
		}
		if err := db.RollbackMigration(databaseURL, steps); err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		fmt.Printf("Rolled back %d migration(s) successfully\n", steps)

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Usage: go run migrate.go force <version>")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		if err := db.ForceMigrationVersion(databaseURL, version); err != nil {
			log.Fatalf("Failed to force migration version: %v", err)
		}
		fmt.Printf("Forced migration version to %d and cleared dirty flag\n", version)

	default:
		log.Fatal("Unknown command. Use 'up', 'down', or 'force'")
	}
}
