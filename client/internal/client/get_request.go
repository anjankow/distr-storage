package client

import (
	"client/internal/config"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type getRsp struct {
	NodeName string          `json:"node"`
	Value    json.RawMessage `json:"value"`
}

func getRequest(key string) (getRsp, error) {

	url := "http://" + config.GetApiAddr() + "/doc?key=" + key

	resp, err := http.Get(url)
	if err != nil {
		return getRsp{}, errors.New("get request failed: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return getRsp{}, errors.New("get return code: " + resp.Status)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return getRsp{}, errors.New("failed to read the response: " + err.Error())
	}

	var responseBody getRsp
	if err := json.Unmarshal(respBytes, &responseBody); err != nil {
		return getRsp{}, errors.New("failed to unmarshal the response: " + err.Error())
	}

	return responseBody, nil
}
