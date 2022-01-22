package app

import (
	nodeproxy "data-access-api/node_proxy"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

type App struct {
	Logger *zap.Logger

	nodes      []string
	collection string
}

func NewApp(l *zap.Logger) (app App, err error) {

	app = App{
		Logger: l,
	}
	return
}

func (a *App) Configure(collection string, nodes []string) {
	// race condition, should be guarded with mutex, but for this use case I'll keep it this way
	a.collection = collection
	a.nodes = nodes
}

func (a App) Insert(key string, value json.RawMessage) (string, time.Time, error) {
	// forward to a node
	node := "node0"
	n := nodeproxy.NodeProxy{
		HostAddr: node,
		Logger:   a.Logger,
	}

	ts, err := n.Insert(key, value)
	if err != nil {
		return "", time.Time{}, err
	}

	return node, ts, nil

}
