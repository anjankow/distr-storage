package app

import (
	nodeproxy "data-access-api/node_proxy"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type App struct {
	Logger *zap.Logger

	collection string
	nodes      []nodeproxy.NodeProxy
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

	for _, nodeAddr := range nodes {
		node := nodeproxy.NodeProxy{
			HostAddr: nodeAddr,
			Logger:   a.Logger,
		}
		a.Logger.Info(fmt.Sprintf(nodeAddr, " waiting..."))
		node.WaitReady()
		a.Logger.Info(fmt.Sprintf(nodeAddr, " ready"))
		a.nodes = append(a.nodes, node)
	}
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
