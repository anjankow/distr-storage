package client

import (
	"client/internal/config"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func getAllRequest() (json.RawMessage, error) {

	url := "http://" + config.GetApiAddr() + "/doc?all=true"

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("get request failed: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("get return code: " + resp.Status)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read the response: " + err.Error())
	}

	var responseBody json.RawMessage
	if err := json.Unmarshal(respBytes, &responseBody); err != nil {
		return nil, errors.New("failed to unmarshal the response: " + err.Error())
	}

	return responseBody, nil
}
