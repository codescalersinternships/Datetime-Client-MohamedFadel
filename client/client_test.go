package client

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSetupEnvironment(t *testing.T) {
	testURL := "http://testserver"
	testPort := "8080"
	testPortGin := "9090"

	SetupEnv(testURL, testPort, testPortGin)

	if os.Getenv("SERVER_URL") != testURL {
		t.Errorf("SERVER_URL not set correctly. Expected %s, got %s", testURL, os.Getenv("SERVER_URL"))
	}

	if os.Getenv("SERVER_PORT") != testPort {
		t.Errorf("SERVER_PORT not set correctly. Expected %s, got %s", testPort, os.Getenv("SERVER_PORT"))
	}

	if os.Getenv("SERVER_PORT_GIN") != testPortGin {
		t.Errorf("SERVER_PORT_GIN not set correctly. Expected %s, got %s", testPortGin, os.Getenv("SERVER_PORT_GIN"))
	}
}

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

	serverURLPlain := serverPlain.URL[:strings.LastIndex(serverPlain.URL, ":")]
	serverPortPlain := serverPlain.URL[strings.LastIndex(serverPlain.URL, ":")+1:]

	serverURLJSON := serverJSON.URL[:strings.LastIndex(serverJSON.URL, ":")]
	serverPortJSON := serverJSON.URL[strings.LastIndex(serverJSON.URL, ":")+1:]

	t.Run("Success - Plain text response", func(t *testing.T) {
		SetupEnv(serverURLPlain, serverPortPlain, serverPortPlain)
		dateTime, err := GetDateTime("other", "text/plain")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if dateTime != "2024-09-19T12:34:56Z" {
			t.Errorf("Expected datetime to be '2024-09-19T12:34:56Z', got %s", dateTime)
		}
	})

	t.Run("Success - JSON response", func(t *testing.T) {
		SetupEnv(serverURLJSON, serverPortJSON, serverPortJSON)
		dateTime, err := GetDateTime("gin", "application/json")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if dateTime != "2024-09-19T12:34:56Z" {
			t.Errorf("Expected datetime to be '2024-09-19T12:34:56Z', got %s", dateTime)
		}
	})

	t.Run("Error - Missing SERVER_URL", func(t *testing.T) {
		os.Setenv("SERVER_URL", "")
		_, err := GetDateTime("gin", "application/json")
		if err == nil {
			t.Fatal("Expected error for missing SERVER_URL, got none")
		}
		os.Setenv("SERVER_URL", serverURLJSON)
	})

	t.Run("Error - Invalid Content-Type", func(t *testing.T) {
		SetupEnv(serverURLPlain, serverPortPlain, serverPortPlain)
		_, err := GetDateTime("other", "unsupported/type")
		if err == nil {
			t.Fatal("Expected error for unsupported content type, got none")
		}
	})

	t.Run("Error - HTTP request failure", func(t *testing.T) {
		SetupEnv("", "invalidPort", "invalidPort")
		_, err := GetDateTime("other", "text/plain")
		if !errors.Is(err, ErrURLandPortMustBeSet) {
			t.Fatalf("Expected ErrURLandPortMustBeSet, got %v", err)
		}
	})

	t.Run("Error - Failed server response", func(t *testing.T) {
		serverURLFail := serverFail.URL[:strings.LastIndex(serverFail.URL, ":")]
		serverPortFail := serverFail.URL[strings.LastIndex(serverFail.URL, ":")+1:]
		SetupEnv(serverURLFail, serverPortFail, serverPortFail)
		_, err := GetDateTime("gin", "application/json")
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

		serverURL := server.URL[:strings.LastIndex(server.URL, ":")]
		serverPort := server.URL[strings.LastIndex(server.URL, ":")+1:]

		SetupEnv(serverURL, serverPort, serverPort)
		dateTime, err := GetDateTime("other", "application/json")
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

		serverURL := server.URL[:strings.LastIndex(server.URL, ":")]
		serverPort := server.URL[strings.LastIndex(server.URL, ":")+1:]

		SetupEnv(serverURL, serverPort, serverPort)
		start := time.Now()
		_, err := GetDateTime("other", "application/json")
		duration := time.Since(start)

		if err == nil {
			t.Fatal("Expected error after max retries, got none")
		}

		if duration < 20*time.Second || duration > 40*time.Second {
			t.Errorf("Expected retry duration to be approximately 30 seconds, got %v", duration)
		}
	})
}
