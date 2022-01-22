package config

import (
	"os"
	"strconv"
)

const (
	defaultNumberOfNodes int = 3
)

var (
	inputFile string
)

func GetInputFilePath() string {
	if inputFile != "" {
		return inputFile
	}

	file := os.Getenv("INPUT_PATH")
	if file != "" {
		inputFile = file
		return inputFile
	}

	return ""
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
