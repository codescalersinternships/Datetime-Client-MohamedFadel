package main

import (
	"os"
	"testing"
)

func TestSetupEnvironment(t *testing.T) {
	// Test values
	testURL := "http://testserver"
	testPort := "8080"
	testPortGin := "9090"

	// Call the function to set up the environment
	SetupEnv(testURL, testPort, testPortGin)

	// Check if the environment variables are set correctly
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
