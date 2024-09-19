package main

import "os"

func SetupEnv(serverURL, serverPort, serverPortGin string) {
	os.Setenv("SERVER_URL", serverURL)
	os.Setenv("SERVER_PORT", serverPort)
	os.Setenv("SERVER_PORT_GIN", serverPortGin)
}

func main() {

}
