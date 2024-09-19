package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("%w", ErrCreatingRequest)
	}

	req.Header.Set("Accept", contentType)

	if contentType != "application/json" && contentType != "text/plain" {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedContentType, contentType)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w", ErrMakingRequest)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%w", ErrReadingResponseBody)
	}

	if contentType == "application/json" {
		var dateTimeResp map[string]string
		err = json.Unmarshal(body, &dateTimeResp)
		if err != nil {
			return "", fmt.Errorf("%w", ErrUnmarshallingJSON)
		}
		return dateTimeResp["datetime"], nil
	}

	return string(body), nil

}
