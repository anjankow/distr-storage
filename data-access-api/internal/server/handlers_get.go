package server

import (
	"data-access-api/internal/app"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

func get(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {

	all := r.URL.Query().Get("all")
	if all == "true" {
		return getAll(a, w, r)
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		return http.StatusBadRequest, errors.New("missing query parameter: key")
	}

	a.Logger.Debug("handling get request", zap.String("key", key))

	node, value, err := a.Get(key)
	if err != nil {
		return http.StatusInternalServerError, errors.New("KEY[" + key + "] get handler failed: " + err.Error())
	}

	var getRsp struct {
		NodeName string          `json:"node"`
		Value    json.RawMessage `json:"value"`
	}
	getRsp.NodeName = node
	getRsp.Value = value

	marshalledRsp, err := json.Marshal(getRsp)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to marshal response body: " + err.Error())
	}

	w.Write(marshalledRsp)

	return http.StatusOK, nil
}

func getAll(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {

	a.Logger.Debug("handling get all request")

	nodeData, err := a.GetAll()
	if err != nil {
		return http.StatusInternalServerError, errors.New("get all handler failed: " + err.Error())
	}

	w.Write(nodeData)

	return http.StatusOK, nil
}
