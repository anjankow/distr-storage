package app

import (
	"go.uber.org/zap"
)

type App struct {
	Logger *zap.Logger
}

func NewApp(l *zap.Logger) (app App, err error) {

	app = App{
		Logger: l,
	}
	return
}
