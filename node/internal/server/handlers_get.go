package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"node/internal/app"

	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func get(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {

	all := r.URL.Query().Get("all")
	if all == "true" {
		return getAll(a, w, r)
	}

	var id string
	var collection string
	if err := multierr.Combine(getDocumentID(r, &id), getCollection(r, &collection)); err != nil {
		return http.StatusBadRequest, err
	}

	a.Logger.Debug("handling get request", zap.String("id", id), zap.String("collection", collection))

	result, err := a.Get(r.Context(), collection, id)
	if err != nil {
		return http.StatusInternalServerError, errors.New("id[" + id + "] get handler failed: " + err.Error())
	}

	if result == nil {
		a.Logger.Debug("writing status: not found", zap.String("id", id))
		return http.StatusNotFound, errors.New("id[" + id + "] document not found")
	}

	w.Write(result)

	return http.StatusOK, nil
}

func getAll(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {

	var collection string
	if err := getCollection(r, &collection); err != nil {
		return http.StatusBadRequest, err
	}

	a.Logger.Debug("handling get all request", zap.String("collection", collection))

	result, err := a.GetAll(r.Context(), collection)
	if err != nil {
		return http.StatusInternalServerError, errors.New("get all handler failed: " + err.Error())
	}

	marshalled, err := json.Marshal(result)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to marshal the response: " + err.Error())
	}

	w.Write(marshalled)

	return http.StatusOK, nil
}
