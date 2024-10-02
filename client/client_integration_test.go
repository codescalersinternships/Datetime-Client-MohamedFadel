package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIntegrationGetDateTime(t *testing.T) {
	tests := []struct {
		name               string
		contentType        string
		serverResponse     interface{}
		expectedStatusCode int
		expectedResult     string
		expectedError      bool
	}{
		{
			name:               "Successful JSON response",
			contentType:        "application/json",
			serverResponse:     map[string]string{"datetime": "2024-09-20T14:30:00Z"},
			expectedStatusCode: http.StatusOK,
			expectedResult:     "2024-09-20T14:30:00Z",
			expectedError:      false,
		},
		{
			name:               "Successful plain text response",
			contentType:        "text/plain",
			serverResponse:     "2024-09-20T14:30:00Z",
			expectedStatusCode: http.StatusOK,
			expectedResult:     "2024-09-20T14:30:00Z",
			expectedError:      false,
		},
		{
			name:               "Server error",
			contentType:        "application/json",
			serverResponse:     nil,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResult:     "",
			expectedError:      true,
		},
		{
			name:               "Invalid JSON response",
			contentType:        "application/json",
			serverResponse:     "invalid json",
			expectedStatusCode: http.StatusOK,
			expectedResult:     "",
			expectedError:      true,
		},
		{
			name:               "Timeout",
			contentType:        "application/json",
			serverResponse:     nil,
			expectedStatusCode: http.StatusOK,
			expectedResult:     "",
			expectedError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tt.contentType)
				w.WriteHeader(tt.expectedStatusCode)

				if tt.name == "Timeout" {
					time.Sleep(11 * time.Second)
					return
				}

				if tt.serverResponse != nil {
					if tt.contentType == "application/json" {
						err := json.NewEncoder(w).Encode(tt.serverResponse)
						if err != nil {
							t.Errorf("Error encoding JSON: %v", err)
						}
					} else {
						_, err := w.Write([]byte(tt.serverResponse.(string)))
						if err != nil {
							t.Errorf("Error writing response: %v", err)
						}
					}
				}
			}))
			defer server.Close()

			client := NewDateTimeClient(server.URL)
			result, err := client.GetDateTime(tt.contentType)

			if tt.expectedError && err == nil {
				t.Errorf("Expected an error, but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
			if result != tt.expectedResult {
				t.Errorf("Expected result %s, but got %s", tt.expectedResult, result)
			}
		})
	}
}
