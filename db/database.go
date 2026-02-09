package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	password := os.Getenv("DB_PASSWORD")
	connectionString := fmt.Sprintf("host=127.0.0.1 port=5432 user=postgres password=%s dbname=car_db sslmode=disable", password)

	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}
	fmt.Println("Successfully connected to the database!")
}
