package config

import (
	"os"
	"strconv"
)

const (
	defaultNumberOfNodes  int    = 3
	defaultCollectionName string = "myCollection"
)

var (
	inputFile      string
	apiAddress     string
	collectionName string
)

func GetInputFilePath() string {

	if inputFile != "" {
		return inputFile
	}

	inputFile := os.Getenv("INPUT_PATH")

	return inputFile
}

func GetNumberOfNodes() int {

	value := os.Getenv("NUMBER_OF_NODES")
	if value == "" {
		return defaultNumberOfNodes
	}

	number, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return defaultNumberOfNodes
	}

	return int(number)

}

func GetApiAddr() string {
	if apiAddress != "" {
		return apiAddress
	}

	apiAddress := os.Getenv("API_ADDR")

	return apiAddress
}

func GetCollectionName() string {
	if collectionName != "" {
		return collectionName
	}

	collectionName := os.Getenv("COLLECTION_NAME")
	if collectionName == "" {
		collectionName = defaultCollectionName
	}

	return collectionName
}
