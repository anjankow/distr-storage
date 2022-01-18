package main

import (
	"log"
	"node/internal/app"
	"node/internal/server"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	logger, err := getLogger()
	if err != nil {
		log.Fatalln("setting up the logger failed: ", err)
		return
	}
	defer logger.Sync()

	logger.Info("service started")

	service, err := app.NewApp(logger)
	if err != nil {
		logger.Fatal("service creation failed: " + err.Error())
		return
	}

	// HTTP SERVER
	ser := server.NewServer(logger, &service)
	if err != nil {
		logger.Fatal("server creation failed: ", zap.Error(err))
	}

	err = ser.Run()

	logger.Info("service finished", zap.Error(err))

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
