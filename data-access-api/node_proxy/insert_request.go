package nodeproxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	maxTrials int = 3
)

func (n NodeProxy) Insert(collection string, id string, content json.RawMessage) (time.Time, error) {
	// passes the object to the given node
	url := "http://" + n.HostAddr + "/doc"

	var body struct {
		Collection string          `json:"collection"`
		ID         string          `json:"id"`
		Content    json.RawMessage `json:"content"`
	}
	body.ID = id
	body.Content = content
	body.Collection = collection

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

		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			n.Logger.Error("failed to read the insert response", zap.Error(err))
			// this error doesn't mean that the insert failed, just the response info can't be read
			err = nil

			break
		}

		var response struct {
			Timestamp time.Time `json:"ts"`
		}
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
