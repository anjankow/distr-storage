package app

import (
	nodeproxy "data-access-api/node_proxy"
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
	n := nodeproxy.NodeProxy{
		HostName: "node0",
		Logger:   a.Logger,
	}

	return n.Insert(key, value)

}
