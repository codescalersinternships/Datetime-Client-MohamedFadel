package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codescalersinternships/Datetime-Client-MohamedFadel/client"
)

func main() {
	serverURL := os.Getenv("SERVER_URL")
	serverPort := os.Getenv("SERVER_PORT")
	serverPortGin := os.Getenv("SERVER_PORT_GIN")

	if serverURL == "" || serverPort == "" || serverPortGin == "" {
		log.Fatal("SERVER_URL, SERVER_PORT, and SERVER_PORT_GIN must be set")
	}

	client.SetupEnv(serverURL, serverPort, serverPortGin)

	serverTypes := []string{"standard", "gin"}
	contentTypes := []string{"application/json", "text/plain"}

	for _, serverType := range serverTypes {
		for _, contentType := range contentTypes {
			dateTime, err := client.GetDateTime(serverType, contentType)
			if err != nil {
				log.Printf("Error fetching datetime from %s with %s: %v\n", serverType, contentType, err)
				continue
			}
			fmt.Printf("Response from %s server with %s: %s\n", serverType, contentType, dateTime)
		}
	}
}
