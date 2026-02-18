package main

import (
	"car-service/db"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scripts/make_admin.go <email>")
		return
	}

	email := os.Args[1]

	// Initialize DB connection
	_, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Update user role to admin
	res, err := db.DB.Exec("UPDATE users SET role='admin' WHERE email=$1", email)
	if err != nil {
		log.Fatalf("Failed to update user role: %v", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		fmt.Printf("User with email '%s' not found.\n", email)
	} else {
		fmt.Printf("Successfully promoted '%s' to ADMIN.\n", email)
	}
}
