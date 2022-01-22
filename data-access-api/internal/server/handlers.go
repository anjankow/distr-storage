package server

import (
	"data-access-api/internal/app"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("all good here"))
}

func insert(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {

	bodyBytes, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return http.StatusInternalServerError, errors.New("can't read the request body: " + err.Error())
	}

	var body struct {
		Key   string          `json:"key"`
		Value json.RawMessage `json:"value"`
	}

	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return http.StatusBadRequest, errors.New("invalid body: " + err.Error())
	}

	a.Logger.Debug("handling insert request", zap.String("key", body.Key), zap.String("value", string(body.Value)))

	node, ts, err := a.Insert(body.Key, body.Value)
	if err != nil {
		return http.StatusInternalServerError, errors.New("KEY[" + body.Key + "] insert handler failed: " + err.Error())
	}

	var insertRsp struct {
		NodeName  string    `json:"node"`
		Timestamp time.Time `json:"time"`
	}
	insertRsp.NodeName = node
	insertRsp.Timestamp = ts

	marshalledRsp, err := json.Marshal(insertRsp)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to marshal response body: " + err.Error())
	}

	w.Write(marshalledRsp)

	return 0, nil
}
