package main

import (
	"car-service/db"
	"fmt"
	"log"
)

func main() {
	_, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	rows, err := db.DB.Query("SELECT id, email, role FROM users")
	if err != nil {
		log.Fatalf("Failed to query users: %v", err)
	}
	defer rows.Close()

	fmt.Println("--- Users in DB ---")
	count := 0
	for rows.Next() {
		var id int
		var email, role string
		if err := rows.Scan(&id, &email, &role); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d | Email: %s | Role: %s\n", id, email, role)
		count++
	}
	if count == 0 {
		fmt.Println("No users found in the database.")
	}
	fmt.Println("-------------------")
}
