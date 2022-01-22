package client

import "go.uber.org/zap"

type Client struct {
	logger *zap.Logger
}

func NewClient(logger *zap.Logger) Client {
	return Client{
		logger: logger,
	}
}

func (c Client) Insert(key string, value interface{}) (node string, err error) {
	return "node0", nil
}

func (c Client) Delete(key string) (node string, err error) {
	return "node0", nil
}

func (c Client) Get(key string) (node string, err error) {
	return "node0", nil
}

func (c Client) GetRange() (node string, err error) {
	return "node0", nil
}
