package config

import "os"

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
