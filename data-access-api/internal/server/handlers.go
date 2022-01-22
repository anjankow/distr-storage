package server

import (
	"data-access-api/internal/app"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("all good here"))
}

func configure(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {

	bodyBytes, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return http.StatusInternalServerError, errors.New("can't read the request body: " + err.Error())
	}

	var body struct {
		Collection string   `json:"collection"`
		Nodes      []string `json:"nodes"`
	}

	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return http.StatusBadRequest, errors.New("invalid body: " + err.Error())
	}

	a.Configure(body.Collection, body.Nodes)

	a.Logger.Debug("updated configuration", zap.String("collection", body.Collection), zap.Any("nodes", body.Nodes))

	return 0, nil
}
