package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Jack-R-Long/alfred/cmd/database"
)

func TestMain(m *testing.M) {
	// Set up test database
	database.Init()
	
	// Run tests
	code := m.Run()
	
	// Clean up
	database.Close()
	os.Exit(code)
}

func TestHealthEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{
			Status: "success",
		}
		response.Data.Message = "API is healthy"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response HealthResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}

	if response.Data.Message != "API is healthy" {
		t.Errorf("Expected message 'API is healthy', got '%s'", response.Data.Message)
	}
}

func TestCreateUserEndpoint(t *testing.T) {
	// Clean up any existing test users
	database.DB.Exec("DELETE FROM users WHERE username LIKE 'testuser%' OR email LIKE 'test%@example.com'")
	
	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		description    string
	}{
		{
			name: "valid user creation",
			requestBody: map[string]string{
				"username": "testuser_create",
				"email":    "testcreate@example.com",
				"password": "testpassword",
			},
			expectedStatus: http.StatusOK,
			description:    "should create user successfully",
		},
		{
			name: "missing username",
			requestBody: map[string]string{
				"email":    "testmissing@example.com",
				"password": "testpassword",
			},
			expectedStatus: http.StatusBadRequest,
			description:    "should fail when username is missing",
		},
		{
			name: "missing email",
			requestBody: map[string]string{
				"username": "testuser_missing",
				"password": "testpassword",
			},
			expectedStatus: http.StatusBadRequest,
			description:    "should fail when email is missing",
		},
		{
			name: "missing password",
			requestBody: map[string]string{
				"username": "testuser_nopass",
				"email":    "testnopass@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			description:    "should fail when password is missing",
		},
		{
			name:           "invalid json",
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
			description:    "should fail with invalid JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			if tt.requestBody != nil {
				json.NewEncoder(&body).Encode(tt.requestBody)
			} else {
				body.WriteString("invalid json")
			}

			req, err := http.NewRequest("POST", "/users", &body)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createUserHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("%s: handler returned wrong status code: got %v want %v",
					tt.description, status, tt.expectedStatus)
			}
		})
	}
}

func TestCreateUserMethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createUserHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}