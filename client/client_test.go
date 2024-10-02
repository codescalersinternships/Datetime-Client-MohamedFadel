package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func mockServer(t *testing.T, contentType string, responseBody string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(statusCode)
		_, err := io.WriteString(w, responseBody)
		if err != nil {
			t.Errorf("Error writing mock response: %v", err)
		}
	}))
}

func TestGetDateTime(t *testing.T) {
	serverPlain := mockServer(t, "text/plain", "2024-09-19T12:34:56Z", http.StatusOK)
	defer serverPlain.Close()

	serverJSON := mockServer(t, "application/json", `{"datetime": "2024-09-19T12:34:56Z"}`, http.StatusOK)
	defer serverJSON.Close()

	serverFail := mockServer(t, "application/json", ``, http.StatusInternalServerError)
	defer serverFail.Close()

	t.Run("Success - Plain text response", func(t *testing.T) {
		client := NewDateTimeClient(serverPlain.URL)
		dateTime, err := client.GetDateTime("text/plain")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if dateTime != "2024-09-19T12:34:56Z" {
			t.Errorf("Expected datetime to be '2024-09-19T12:34:56Z', got %s", dateTime)
		}
	})

	t.Run("Success - JSON response", func(t *testing.T) {
		client := NewDateTimeClient(serverJSON.URL)
		dateTime, err := client.GetDateTime("application/json")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if dateTime != "2024-09-19T12:34:56Z" {
			t.Errorf("Expected datetime to be '2024-09-19T12:34:56Z', got %s", dateTime)
		}
	})

	t.Run("Error - Invalid Content-Type", func(t *testing.T) {
		client := NewDateTimeClient(serverPlain.URL)
		_, err := client.GetDateTime("unsupported/type")
		if err == nil {
			t.Fatal("Expected error for unsupported content type, got none")
		}
	})

	t.Run("Error - HTTP request failure", func(t *testing.T) {
		client := NewDateTimeClient("http://invalidurl:1234")
		_, err := client.GetDateTime("text/plain")
		if err == nil {
			t.Fatal("Expected error for invalid URL, got none")
		}
	})

	t.Run("Error - Failed server response", func(t *testing.T) {
		client := NewDateTimeClient(serverFail.URL)
		_, err := client.GetDateTime("application/json")
		if err == nil {
			t.Fatal("Expected error for failed server response, got none")
		}
	})

	t.Run("Retry mechanism - Success after failures", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			if attempts <= 2 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, err := io.WriteString(w, `{"datetime": "2024-09-19T12:34:56Z"}`)
			if err != nil {
				t.Errorf("Error writing retry response: %v", err)
			}
		}))
		defer server.Close()

		client := NewDateTimeClient(server.URL)
		dateTime, err := client.GetDateTime("application/json")
		if err != nil {
			t.Fatalf("Expected no error after retries, got %v", err)
		}
		if dateTime != "2024-09-19T12:34:56Z" {
			t.Errorf("Expected datetime to be '2024-09-19T12:34:56Z', got %s", dateTime)
		}
		if attempts != 3 {
			t.Errorf("Expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("Retry mechanism - Failure after max retries", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewDateTimeClient(server.URL)
		start := time.Now()
		_, err := client.GetDateTime("application/json")
		duration := time.Since(start)

		if err == nil {
			t.Fatal("Expected error after max retries, got none")
		}

		if duration < 20*time.Second || duration > 40*time.Second {
			t.Errorf("Expected retry duration to be approximately 30 seconds, got %v", duration)
		}
	})
}
