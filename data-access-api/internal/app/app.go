package app

import (
	nodeproxy "data-access-api/node_proxy"
	"encoding/json"
	"errors"
	"fmt"
	"hash/adler32"
	"sync"
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

	ts, err := node.Insert(a.collection, key, value)
	if err != nil {
		return "", time.Time{}, err
	}

	return node.HostAddr, ts, nil

}

func (a App) Get(key string) (string, json.RawMessage, error) {

	nodeIdx, err := a.getNodeIdx(key)
	if err != nil {
		return "", nil, err
	}

	node := a.nodes[nodeIdx]

	value, err := node.Get(a.collection, key)
	if err != nil {
		return "", nil, err
	}

	return node.HostAddr, value, nil

}

func (a App) GetAll() ([]byte, error) {

	var wg sync.WaitGroup
	var mutex sync.Mutex

	type nodeData struct {
		NodeName string
		Data     json.RawMessage
	}
	var combined struct {
		Data []nodeData `json:"node_data"`
	}
	combined.Data = make([]nodeData, len(a.nodes))

	for i := 0; i < len(a.nodes); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			nodeIdx := i

			node := a.nodes[nodeIdx]

			data, err := node.GetAll(a.collection)
			if err != nil {
				a.Logger.Error("failed to getall from "+node.HostAddr, zap.Error(err))
			}

			mutex.Lock()
			defer mutex.Unlock()

			combined.Data[i].NodeName = node.HostAddr
			combined.Data[i].Data = data
		}(i)
	}
	wg.Wait()

	combinedBin, err := json.Marshal(combined)
	if err != nil {
		return nil, errors.New("failed to marshal combined node response: " + err.Error())
	}

	return combinedBin, nil
}
