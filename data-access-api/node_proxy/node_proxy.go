package nodeproxy

import (
	"bytes"
	"data-access-api/internal/config"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	maxTrials int = 3
)

type NodeProxy struct {
	HostName string
	Logger   *zap.Logger
}

type insertResponse struct {
	Timestamp time.Time `json:"ts"`
}

func (n NodeProxy) Insert(id string, content json.RawMessage) (time.Time, error) {
	// passes the object to the given node
	url := "http://" + n.HostName + config.NodePort + "/doc"

	var body struct {
		Collection string          `json:"collection"`
		ID         string          `json:"id"`
		Content    json.RawMessage `json:"content"`
	}
	body.ID = id
	body.Content = content
	body.Collection = config.GetCollectionName()

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return time.Time{}, errors.New("failed to marshal body when sending the request: " + err.Error())
	}
	reqBody := bytes.NewBuffer(marshalledBody)

	req, err := http.NewRequest(http.MethodPut, url, reqBody)
	if err != nil {
		return time.Time{}, errors.New("NodeProxy: failed to create insert request: " + err.Error())
	}

	var insertTime time.Time

	for i := 0; i < maxTrials; i++ {
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			n.Logger.Debug("failed to do request: "+err.Error(), zap.Int("trial", i))
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			n.Logger.Debug("insert request return code: "+resp.Status, zap.Int("trial", i))
			continue
		}

		var respBytes []byte
		n.Logger.Error("response ", zap.Any("cont len", resp.ContentLength))

		if _, err = resp.Body.Read(respBytes); err != nil {
			n.Logger.Error("failed to read the insert response", zap.Error(err))
			// this error doesn't mean that the insert failed, just the response info can't be read
			err = nil

			break
		}
		n.Logger.Info("reposne bytes: " + string(respBytes))

		// var response struct {
		// 	Timestamp time.Time `json:"ts"`
		// }
		response := insertResponse{}
		if err := json.Unmarshal(respBytes, &response); err != nil {
			n.Logger.Error("failed to unmarshal the insert response", zap.Error(err), zap.Any("response", respBytes))
			// this error doesn't mean that the insert failed, just the response info can't be read
			err = nil
		}
		insertTime = response.Timestamp

		break
	}

	if err != nil {
		return time.Time{}, errors.New("NodeProxy: " + err.Error())
	}

	return insertTime, err
}
