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

	insertedTime, err := a.Insert(r.Context(), body.Collection, body.ID, body.Content)
	if err != nil {
		return http.StatusInternalServerError, errors.New("id[" + body.ID + "] insert handler failed: " + err.Error())
	}

	marshalledTime, err := insertedTime.MarshalText()
	if err != nil {
		a.Logger.Warn("failed to marshal time: "+err.Error(), zap.String("id", body.ID))
	}
	w.Write(marshalledTime)

	return 0, nil
}

func get(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {

	var id string
	var collection string
	if err := multierr.Combine(getDocumentID(r, &id), getCollection(r, &collection)); err != nil {
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

func delete(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {

	var id string
	var collection string
	if err := multierr.Combine(getDocumentID(r, &id), getCollection(r, &collection)); err != nil {
		return http.StatusBadRequest, err
	}

	a.Logger.Debug("handling delete request", zap.String("id", id))

	err := a.Delete(r.Context(), collection, id)
	if err != nil {
		return http.StatusInternalServerError, errors.New("id[" + id + "] delete handler failed: " + err.Error())
	}

	return 0, nil
}

func getCollection(r *http.Request, collection *string) error {
	coll := r.URL.Query().Get("collection")

	if coll == "" {
		return errors.New("missing query parameter: collection")
	}

	*collection = coll
	return nil
}

func getDocumentID(r *http.Request, docID *string) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		return errors.New("missing query parameter: id")
	}
	*docID = id
	return nil
}
