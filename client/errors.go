package client

import "errors"

var (
	ErrURLandPortMustBeSet = errors.New("SERVER_URL and SERVER_PORT must be set in environment variables")

	ErrCreatingRequest = errors.New("error creating request")

	ErrMakingRequest = errors.New("error making request")

	ErrReadingResponseBody = errors.New("error reading response body")

	ErrUnmarshallingJSON = errors.New("error unmarshalling JSON")
)
