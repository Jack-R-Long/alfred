package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jack-R-Long/alfred/cmd/database"
)

func setupTestUser(t *testing.T, username, email, password string) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		t.Fatal("Failed to hash password:", err)
	}

	_, err = database.DB.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", 
		username, email, hashedPassword)
	if err != nil {
		t.Fatal("Failed to create test user:", err)
	}
}

func cleanupTestUser(username string) {
	database.DB.Exec("DELETE FROM users WHERE username = ?", username)
}

func TestUserHandlerGet(t *testing.T) {
	// Setup test user
	testUsername := "gettest_unique"
	testEmail := "gettest_unique@example.com"
	cleanupTestUser(testUsername) // Clean up any existing user
	setupTestUser(t, testUsername, testEmail, "password123")
	defer cleanupTestUser(testUsername)

	tests := []struct {
		name           string
		username       string
		expectedStatus int
		description    string
	}{
		{
			name:           "existing user",
			username:       testUsername,
			expectedStatus: http.StatusOK,
			description:    "should return user details for existing user",
		},
		{
			name:           "non-existing user",
			username:       "nonexistent",
			expectedStatus: http.StatusNotFound,
			description:    "should return 404 for non-existing user",
		},
		{
			name:           "empty username",
			username:       "",
			expectedStatus: http.StatusBadRequest,
			description:    "should return 400 for empty username",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/users/"+tt.username, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(userHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("%s: handler returned wrong status code: got %v want %v",
					tt.description, status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response HealthResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatal("Failed to decode response:", err)
				}

				if response.Status != "success" {
					t.Errorf("Expected status 'success', got '%s'", response.Status)
				}
			}
		})
	}
}

func TestUserHandlerPut(t *testing.T) {
	// Setup test user
	testUsername := "puttest_unique"
	testEmail := "puttest_unique@example.com"
	cleanupTestUser(testUsername) // Clean up any existing user
	setupTestUser(t, testUsername, testEmail, "oldpassword")
	defer cleanupTestUser(testUsername)

	tests := []struct {
		name           string
		username       string
		requestBody    map[string]string
		expectedStatus int
		description    string
	}{
		{
			name:     "valid update",
			username: testUsername,
			requestBody: map[string]string{
				"email":    "newemail@example.com",
				"password": "newpassword123",
			},
			expectedStatus: http.StatusOK,
			description:    "should update user successfully",
		},
		{
			name:     "missing email",
			username: testUsername,
			requestBody: map[string]string{
				"password": "newpassword123",
			},
			expectedStatus: http.StatusBadRequest,
			description:    "should fail when email is missing",
		},
		{
			name:     "missing password",
			username: testUsername,
			requestBody: map[string]string{
				"email": "newemail@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			description:    "should fail when password is missing",
		},
		{
			name:     "non-existing user",
			username: "nonexistent",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusNotFound,
			description:    "should return 404 for non-existing user",
		},
		{
			name:           "empty username",
			username:       "",
			requestBody:    map[string]string{},
			expectedStatus: http.StatusBadRequest,
			description:    "should return 400 for empty username",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tt.requestBody)

			req, err := http.NewRequest("PUT", "/users/"+tt.username, &body)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(userHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("%s: handler returned wrong status code: got %v want %v",
					tt.description, status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response HealthResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatal("Failed to decode response:", err)
				}

				if response.Status != "success" {
					t.Errorf("Expected status 'success', got '%s'", response.Status)
				}
			}
		})
	}
}

func TestUserHandlerInvalidMethod(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/users/testuser", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	
	hash1, err := hashPassword(password)
	if err != nil {
		t.Fatal("Failed to hash password:", err)
	}

	if hash1 == "" {
		t.Error("Hash should not be empty")
	}

	if hash1 == password {
		t.Error("Hash should be different from original password")
	}

	// Test that same password produces different hashes (due to salt)
	hash2, err := hashPassword(password)
	if err != nil {
		t.Fatal("Failed to hash password:", err)
	}

	if hash1 == hash2 {
		t.Error("Different calls should produce different hashes due to salt")
	}
}

func TestHashPasswordEmpty(t *testing.T) {
	_, err := hashPassword("")
	if err != nil {
		t.Error("Hashing empty password should not error, got:", err)
	}
}