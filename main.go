package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/codescalersinternships/Datetime-Client-MohamedFadel/client"
)

func main() {
	serverURL := flag.String("url", "http://localhost:8000", "Server URL")
	contentType := flag.String("type", "application/json", "Content type (application/json or text/plain)")
	flag.Parse()

	dateTimeClient := client.NewDateTimeClient(*serverURL)

	dateTime, err := dateTimeClient.GetDateTime(*contentType)
	if err != nil {
		log.Fatalf("Error fetching datetime: %v\n", err)
	}

	fmt.Printf("Response from server: %s\n", dateTime)
}
