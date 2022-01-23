package main

import (
	"client/internal/client"
	"client/internal/config"
	"client/internal/tester"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func readInputDataFile() (content map[string]interface{}, err error) {
	path := config.GetInputFilePath()
	if path == "" {
		return nil, errors.New("input file path not configured")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, err
	}

	return rawData, nil
}

func main() {

	logger, err := getLogger()
	if err != nil {
		log.Fatalln("setting up the logger failed: ", err)
		return
	}
	defer logger.Sync()

	logger.Info("test started")

	values, err := readInputDataFile()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	clnt := client.NewClient(logger)

	test := tester.NewTester(logger, &clnt)
	test.Run(values)

	file, err := os.Create("report.md")

	report := test.GenerateReport()
	if err != nil {
		logger.Info("failed to create report file, the report will be displayed on the console")
		fmt.Print(report)
	} else {
		logger.Info("report written to report.md")
		file.WriteString(report)
	}

	logger.Info("test finished", zap.Error(err))

}

func getLogger() (*zap.Logger, error) {
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zap.FatalLevel),
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	config.Development = true
	config.Level.SetLevel(zap.DebugLevel)

	logger, err := config.Build()
	return logger.WithOptions(options...), err
}
