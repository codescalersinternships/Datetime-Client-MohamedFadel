package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
)

func SetupEnv(serverURL, serverPort, serverPortGin string) {
	os.Setenv("SERVER_URL", serverURL)
	os.Setenv("SERVER_PORT", serverPort)
	os.Setenv("SERVER_PORT_GIN", serverPortGin)
}

func GetDateTime(serverType, contentType string) (string, error) {
	serverURL := os.Getenv("SERVER_URL")
	var serverPort string
	if serverType == "gin" {
		serverPort = os.Getenv("SERVER_PORT_GIN")
	} else {
		serverPort = os.Getenv("SERVER_PORT")
	}

	if serverURL == "" || serverPort == "" {
		return "", fmt.Errorf("%w", ErrURLandPortMustBeSet)
	}

	url := fmt.Sprintf("%s:%s/datetime", serverURL, serverPort)

	var result string
	operation := func() error {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("%w", ErrCreatingRequest)
		}

		req.Header.Set("Accept", contentType)

		if contentType != "application/json" && contentType != "text/plain" {
			return fmt.Errorf("%w: %s", ErrUnsupportedContentType, contentType)
		}

		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("%w", ErrMakingRequest)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("%w", ErrReadingResponseBody)
		}

		if contentType == "application/json" {
			var dateTimeResp map[string]string
			err = json.Unmarshal(body, &dateTimeResp)
			if err != nil {
				return fmt.Errorf("%w", ErrUnmarshallingJSON)
			}
			result = dateTimeResp["datetime"]
		} else {
			result = string(body)
		}

		return nil
	}

	backOff := backoff.NewExponentialBackOff()
	backOff.MaxElapsedTime = 30 * time.Second

	err := backoff.Retry(operation, backOff)
	if err != nil {
		return "", err
	}

	return result, nil
}
