package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"car-app/db"
	"car-app/handlers"
)

func main() {

	//

	db.Connect()
	defer db.DB.Close()

	router := mux.NewRouter()
	router.Use(loggingMiddleware)
	router.HandleFunc("/cars", handlers.GetCars).Methods("GET")
	router.HandleFunc("/cars", handlers.CreateCar).Methods("POST")
	router.HandleFunc("/cars/{id}", handlers.GetCar).Methods("GET")
	router.HandleFunc("/cars/{id}", handlers.UpdateCar).Methods("PUT")
	secure := router.PathPrefix("/admin").Subrouter()
	secure.Use(authMiddleware)
	secure.HandleFunc("/delete/{id}", handlers.DeleteCar).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request Started: [%s] %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretToken := os.Getenv("ADMIN_TOKEN")
		if secretToken == "" {
			secretToken = "Infobell"
		}
		userToken := r.Header.Get("X-API-Key")
		if userToken != secretToken {
			http.Error(w, "Forbidden: Invalid API Key", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
