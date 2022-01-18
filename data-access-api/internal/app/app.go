package app

import (
	"data-access-api/node"
	"encoding/json"

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

func (a App) Insert(key string, value json.RawMessage) error {
	// forward to a node
	n := node.NodeProxy{
		HostName: "mongo2",
		Logger:   a.Logger,
	}
	return n.Insert(key, value)

}
