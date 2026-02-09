package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres" // Fallback
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, 5432, "postgres", password, "Cars")

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	// Optimize: Connection Pooling
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	err = DB.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to database with pooling enabled!")
	return DB, nil
}
