package client

import (
	"bytes"
	"client/internal/config"
	"encoding/json"
	"errors"
	"net/http"
)

type configReq struct {
	Collection string   `json:"collection"`
	Nodes      []string `json:"nodes"`
}

func configRequest(body configReq) error {

	url := "http://" + config.GetApiAddr() + "/config"

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return errors.New("failed to marshal body when sending the request: " + err.Error())
	}
	reqBody := bytes.NewBuffer(marshalledBody)

	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return errors.New("failed to create config request: " + err.Error())
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("config return code: " + resp.Status)
	}

	return nil
}
