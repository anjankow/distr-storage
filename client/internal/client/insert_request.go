package client

import (
	"bytes"
	"client/internal/config"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type insertReq struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

type insertRsp struct {
	NodeName  string    `json:"node"`
	Timestamp time.Time `json:"time"`
}

func insertRequest(body insertReq) (insertRsp, error) {

	url := "http://" + config.GetApiAddr() + "/doc"

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return insertRsp{}, errors.New("failed to marshal body when sending the request: " + err.Error())
	}
	reqBody := bytes.NewBuffer(marshalledBody)

	req, err := http.NewRequest(http.MethodPut, url, reqBody)
	if err != nil {
		return insertRsp{}, errors.New("failed to create insert request: " + err.Error())
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return insertRsp{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return insertRsp{}, errors.New("insert return code: " + resp.Status)
	}

	var respBytes []byte
	if _, err = reqBody.Read(respBytes); err != nil {
		return insertRsp{}, errors.New("failed to read the response: " + err.Error())
	}

	fmt.Println(string(respBytes))
	var responseBody insertRsp
	if err := json.Unmarshal(respBytes, &responseBody); err != nil {
		return insertRsp{}, errors.New("failed to unmarshal the response: " + err.Error())
	}

	return responseBody, nil
}
