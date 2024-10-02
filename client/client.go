package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
}

type DateTimeClient struct {
	ServerURL string
}

func NewDateTimeClient(serverURL string) *DateTimeClient {
	return &DateTimeClient{
		ServerURL: serverURL,
	}
}

// GetDateTime retrieves the current date and time from the specified server.
// It supports both standard and Gin server types and can handle JSON or plain text responses.
//
// Parameters:
//   - serverType: The type of server to query ("standard" or "gin")
//   - contentType: The desired content type for the response ("application/json" or "text/plain")
//
// Returns:
//   - string: The retrieved datetime as a string
//   - error: An error if the request fails or the response is invalid
//
// The function uses an exponential backoff retry mechanism for failed requests.
func (c *DateTimeClient) GetDateTime(contentType string) (string, error) {
	if contentType != "application/json" && contentType != "text/plain" {
		logger.Error("Unsupported content type", zap.String("contentType", contentType))
		return "", fmt.Errorf("%w: %s", ErrUnsupportedContentType, contentType)
	}

	logger.Info("Preparing to fetch datetime",
		zap.String("url", c.ServerURL),
		zap.String("contentType", contentType))

	var result string
	operation := func() error {
		req, err := http.NewRequest("GET", c.ServerURL, nil)
		if err != nil {
			logger.Error("Failed to create request", zap.Error(err))
			return fmt.Errorf("%w", ErrCreatingRequest)
		}

		req.Header.Set("Accept", contentType)

		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			logger.Error("Failed to make request", zap.Error(err))
			return fmt.Errorf("%w", ErrMakingRequest)
		}
		defer resp.Body.Close()

		logger.Info("Received response",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("contentType", resp.Header.Get("Content-Type")))

		if resp.StatusCode != http.StatusOK {
			logger.Warn("Unexpected status code", zap.Int("statusCode", resp.StatusCode))
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("Failed to read response body", zap.Error(err))
			return fmt.Errorf("%w", ErrReadingResponseBody)
		}

		if contentType == "application/json" {
			var dateTimeResp map[string]string
			err = json.Unmarshal(body, &dateTimeResp)
			if err != nil {
				logger.Error("Failed to unmarshal JSON response", zap.Error(err))
				return fmt.Errorf("%w", ErrUnmarshallingJSON)
			}
			result = dateTimeResp["datetime"]
		} else {
			result = string(body)
		}

		logger.Info("Successfully retrieved datetime", zap.String("result", result))
		return nil
	}

	backOff := backoff.NewExponentialBackOff()
	backOff.MaxElapsedTime = 30 * time.Second

	err := backoff.Retry(operation, backOff)
	if err != nil {
		logger.Error("All retry attempts failed", zap.Error(err))
		return "", err
	}

	return result, nil
}
