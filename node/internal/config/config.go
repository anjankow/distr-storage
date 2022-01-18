package config

import "os"

const (
	defaultLocalPort = ":8080"

	DefaultDbPort       = ":27017"
	defaultDatabaseName = "myDB"
)

var (
	port          string
	connectionURI string
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

func GetDbConnectionURI() string {
	if connectionURI != "" {
		return connectionURI
	}

	connectionURI = os.Getenv("DB_URI")

	return connectionURI
}

func GetDatabaseName() string {
	return defaultDatabaseName
}
