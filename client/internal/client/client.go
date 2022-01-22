package client

import (
	"encoding/json"
	"errors"

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

func (c Client) Get(key string) (node string, err error) {
	return "node0", nil
}

func (c Client) GetAllData() (data string, err error) {
	return "all the data", nil
}

func (c Client) GetRange() (node string, err error) {
	return "node0", nil
}

func (c Client) ConfigureSystem(numberOfNodes int, collection string) error {

	return nil
}
