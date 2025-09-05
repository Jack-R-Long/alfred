package main

import (
	"database/sql"
	"encoding/json"
	"strings"

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
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract user details from request body
		var user struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if user.Username == "" || user.Email == "" || user.Password == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}
		
		sqlTx, err := database.DB.BeginTx(r.Context(), nil)
		if err != nil {
			http.Error(w, "Failed to create db tx", http.StatusNotFound)
			return
		}
		
		_, err = sqlTx.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", user.Username, user.Email , user.Password)
		
		if err != nil {
			sqlTx.Rollback()
			http.Error(w, "Failed to insert user", http.StatusInternalServerError)
			return
		}

		if err = sqlTx.Commit(); err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{
			Status: "success",
			Data: struct {
				Message string `json:"message"`
			}{
				Message: "User created",
			},
		})
	})

	// Get or update user
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		username := strings.TrimPrefix(r.URL.Path, "/users/")
		
		if username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			// Fetch user details from the database
			var (
				id			  int
				email			  string
			)
			err := database.DB.QueryRow("SELECT id, email FROM users WHERE username = ?", username).Scan(&id, &email)
			
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "User not found", http.StatusNotFound)
				} else {
					http.Error(w, "Failed to query user", http.StatusInternalServerError)
				}
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(HealthResponse{
				Status: "success",
				Data: struct {
					Message  string `json:"message"`
				}{
					Message:  "User found id: " + string(rune(id)) + " with email " + email,
				},
			})
		case http.MethodPut:
			// PUT logic to update user
			var user struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}

			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			result, err := database.DB.Exec("UPDATE users SET email = ?, password_hash = ? WHERE username = ?", user.Email, user.Password, username)
			if err != nil {
				http.Error(w, "Failed to update user", http.StatusInternalServerError)
				return
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				http.Error(w, "Failed to retrieve affected rows", http.StatusInternalServerError)
				return
			}
			if rowsAffected == 0 {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			
			// Respond with success message
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(HealthResponse{
				Status: "success",
				Data: struct {
					Message string `json:"message"`
				}{
					Message: "User updated",
				},
			})
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}