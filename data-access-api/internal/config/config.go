package config

import "os"

const (
	NodePort          = ":8080"
	DefaultCollection = "DeckList"
	DbName            = "myDB"

	defaultLocalPort = ":8082"
)

var (
	port string
)

// GetPort returns port prepended with `:`
func GetPort() string {
	if port != "" {
		return port
	}

	portNum := os.Getenv("PORT")
	if portNum != "" {
		port = ":" + portNum
		return port
	}

	port = defaultLocalPort
	return port
}
