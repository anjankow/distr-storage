package client

import (
	"client/internal/config"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type deleteRsp struct {
	NodeName string `json:"node"`
}

func deleteRequest(key string) (deleteRsp, error) {

	url := "http://" + config.GetApiAddr() + "/doc?key=" + key

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return deleteRsp{}, errors.New("failed to create delete request: " + err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return deleteRsp{}, errors.New("delete request failed: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return deleteRsp{}, errors.New("delete return code: " + resp.Status)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return deleteRsp{}, errors.New("failed to read the response: " + err.Error())
	}

	var responseBody deleteRsp
	if err := json.Unmarshal(respBytes, &responseBody); err != nil {
		return deleteRsp{}, errors.New("failed to unmarshal the response: " + err.Error())
	}

	return responseBody, nil
}
