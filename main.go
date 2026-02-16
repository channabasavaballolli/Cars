package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"car-service/db"
	"car-service/graph"
	"car-service/handlers"
	"car-service/middleware"

	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"

	_ "net/http/pprof" // Import for side-effects
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("Method: %s, URL: %s, Duration: %s", r.Method, r.URL, duration)
	})
}

func main() {
	// Start pprof server on port 6060
	go func() {
		log.Println("Starting pprof server on :6060")
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	_, _ = db.InitDB() // Initialize DB (Skeleton)

	// Reset Database on Startup (As requested)
	// if err := db.ResetDB(); err != nil {
	// 	log.Printf("Warning: Failed to reset DB: %v", err)
	// }

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.HandleFunc("/cars", handlers.GetCars).Methods("GET")
	r.HandleFunc("/cars", handlers.CreateCar).Methods("POST")
	r.HandleFunc("/cars/{id}", handlers.GetCar).Methods("GET")
	r.HandleFunc("/cars/{id}", handlers.UpdateCar).Methods("PUT")

	// Protect DELETE route (Now handled by GraphQL or could be updated here if REST is still used)
	r.HandleFunc("/cars/{id}", handlers.DeleteCar).Methods("DELETE")

	// GraphQL Endpoint
	schema, err := graph.InitSchema()
	if err != nil {
		log.Fatalf("Failed to create GraphQL schema: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	r.Handle("/graphql", middleware.AuthMiddleware(h))

	fmt.Println("Server starting...")
	log.Fatal(http.ListenAndServe(":8000", r))
}
