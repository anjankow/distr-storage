package server

import (
	"data-access-api/internal/app"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

func delete(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {

	key := r.URL.Query().Get("key")
	if key == "" {
		return http.StatusBadRequest, errors.New("missing query parameter: key")
	}

	a.Logger.Debug("handling delete request", zap.String("key", key))

	node, err := a.Delete(key)
	if err != nil {
		return http.StatusInternalServerError, errors.New("KEY[" + key + "] delete handler failed: " + err.Error())
	}

	var deleteRsp struct {
		NodeName string `json:"node"`
	}
	deleteRsp.NodeName = node

	marshalledRsp, err := json.Marshal(deleteRsp)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to marshal response body: " + err.Error())
	}

	w.Write(marshalledRsp)

	return http.StatusOK, nil
}
