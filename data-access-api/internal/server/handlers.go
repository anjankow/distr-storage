package server

import (
	"data-access-api/internal/app"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

const (
	contentType = "application/json"
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

	return 0, nil
}
