package app

import (
	nodeproxy "data-access-api/node_proxy"
	"encoding/json"
	"errors"
	"fmt"
	"hash/adler32"
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
		a.Logger.Info(fmt.Sprint(nodeAddr, " waiting..."))
		node.WaitReady()
		a.Logger.Info(fmt.Sprint(nodeAddr, " ready"))
		a.nodes = append(a.nodes, node)
	}
}

func (a App) getNodeIdx(key string) (int, error) {
	hashFunc := adler32.New()
	if _, err := hashFunc.Write(([]byte)(key)); err != nil {
		return 0, errors.New("failed to initialize hash function: " + err.Error())
	}

	sum := hashFunc.Sum32()
	if len(a.nodes) == 0 {
		return 0, errors.New("no nodes configured")
	}
	nodeIdx := sum % uint32(len(a.nodes))
	a.Logger.Debug("hash function", zap.Uint32("node_idx", nodeIdx), zap.String("key", key), zap.Uint32("hash", sum))

	return int(nodeIdx), nil
}

func (a App) Insert(key string, value json.RawMessage) (string, time.Time, error) {

	nodeIdx, err := a.getNodeIdx(key)
	if err != nil {
		return "", time.Time{}, err
	}

	node := a.nodes[nodeIdx]

	ts, err := node.Insert(key, value)
	if err != nil {
		return "", time.Time{}, err
	}

	return node.HostAddr, ts, nil

}
