package nodeproxy

import (
	"bytes"
	"data-access-api/internal/config"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

const (
	maxTrials int = 3
)

type NodeProxy struct {
	HostName string
	Logger   *zap.Logger
}

func (n NodeProxy) Insert(key string, value json.RawMessage) error {
	// passes the object to the given node
	url := "http://" + n.HostName + config.NodePort + "/insert"

	var body struct {
		Collection string          `json:"collection"`
		Key        string          `json:"key"`
		Value      json.RawMessage `json:"value"`
	}
	body.Key = key
	body.Value = value
	body.Collection = config.DefaultCollection

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return errors.New("failed to marshal body when sending the request: " + err.Error())
	}
	reqBody := bytes.NewBuffer(marshalledBody)

	req, err := http.NewRequest(http.MethodPut, url, reqBody)
	if err != nil {
		return errors.New("NodeProxy: failed to create insert request: " + err.Error())
	}

	var resp *http.Response
	for i := 0; i < maxTrials; i++ {
		resp, err = http.DefaultClient.Do(req)

		if err != nil {
			n.Logger.Debug("failed to do request: "+err.Error(), zap.Int("trial", i))
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = errors.New("return code: " + resp.Status)
			n.Logger.Debug("request return code: "+resp.Status, zap.Int("trial", i))
			continue
		}

		// no errors or failed status code - no need to repeat the request
		break
	}

	if err != nil {
		return errors.New("NodeProxy: " + err.Error())
	}

	return nil
}
