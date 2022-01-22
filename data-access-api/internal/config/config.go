package config

import "os"

const (
	NodePort          = ":8080"
	defaultCollection = "testCollection"
	DbName            = "myDB"

	defaultLocalPort = ":8082"
)

var (
	port           string
	collectionName string
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

func GetCollectionName() string {
	if collectionName != "" {
		return collectionName
	}

	collectionName = defaultCollection
	return collectionName
}
