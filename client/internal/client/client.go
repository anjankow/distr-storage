package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type Client struct {
	logger *zap.Logger
}

func NewClient(logger *zap.Logger) Client {
	return Client{
		logger: logger,
	}
}

func (c Client) Insert(key string, value interface{}) (node string, err error) {

	valueBytes, err := json.Marshal(value)
	if err != nil {
		return "", errors.New("value marshal error: " + err.Error())
	}

	reqBody := insertReq{
		Key:   key,
		Value: json.RawMessage(valueBytes),
	}
	rspBody, err := insertRequest(reqBody)
	if err != nil {
		return "", err
	}

	return rspBody.NodeName, nil
}

func (c Client) Delete(key string) (node string, err error) {
	return "node0", nil
}

func (c Client) Get(key string) (value json.RawMessage, node string, err error) {
	rspBody, err := getRequest(key)
	if err != nil {
		return nil, "", err
	}

	return rspBody.Value, rspBody.NodeName, nil
}

func (c Client) GetAllData() (data string, err error) {
	return "all the data", nil
}

func (c Client) GetRange() (node string, err error) {
	return "node0", nil
}

func (c Client) ConfigureSystem(numberOfNodes int, collection string) error {

	nodes := []string{}
	for i := 0; i < numberOfNodes; i++ {
		nodes = append(nodes, fmt.Sprintf("node%d:8080", i))
	}
	reqBody := configReq{
		Collection: collection,
		Nodes:      nodes,
	}

	if err := configRequest(reqBody); err != nil {
		return err
	}

	c.logger.Info("successfully configured")

	return nil
}

func (c Client) Wait() {
	for {
		err := healthRequest()
		if err != nil {
			continue
		}

		break
	}
	c.logger.Info("system responding")
}
