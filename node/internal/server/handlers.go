package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"node/internal/app"

	"go.uber.org/multierr"
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
		Collection string          `json:"collection"`
		ID         string          `json:"id"`
		Content    json.RawMessage `json:"content"`
	}

	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return http.StatusBadRequest, errors.New("invalid body: " + err.Error())
	}

	a.Logger.Debug("handling insert request", zap.String("id", body.ID), zap.String("content", string(body.Content)))

	if err = a.Insert(r.Context(), body.Collection, body.ID, body.Content); err != nil {
		return http.StatusInternalServerError, errors.New("id[" + body.ID + "] insert handler failed: " + err.Error())
	}

	return 0, nil
}

func get(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {

	id := r.URL.Query().Get("id")
	var idErr error
	if id == "" {
		idErr = errors.New("missing query parameter: id")
	}

	collection := r.URL.Query().Get("collection")
	var collectionErr error
	if collection == "" {
		collectionErr = errors.New("missing query parameter: collection")
	}

	if err := multierr.Combine(idErr, collectionErr); err != nil {
		return http.StatusBadRequest, err
	}

	a.Logger.Debug("handling get request", zap.String("id", id))

	result, err := a.Get(r.Context(), collection, id)
	if err != nil {
		return http.StatusInternalServerError, errors.New("id[" + id + "] get handler failed: " + err.Error())
	}

	if result == nil {
		a.Logger.Debug("writing status: not found", zap.String("id", id))
		return http.StatusNotFound, errors.New("id[" + id + "] document not found")
	}

	w.Write(result)

	return 0, nil
}
