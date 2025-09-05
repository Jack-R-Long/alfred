package main

import (
	"encoding/json"

	// "fmt"
	"log"
	"net/http"

	"github.com/Jack-R-Long/alfred/cmd/database"
)

func main() {
	// Initialize the database
	database.Init()
	defer database.Close()

	// Define a simple health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{
			Status: "success",
		}
		response.Data.Message = "API is healthy"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	})

	// Create a user
	http.HandleFunc("/users", createUserHandler)

	// Get or update user
	http.HandleFunc("/users/", userHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}