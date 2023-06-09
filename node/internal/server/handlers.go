package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"node/internal/app"
	"time"

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

	a.Logger.Debug("handling insert request", zap.String("id", body.ID), zap.String("collection", body.Collection))

	insertedTime, err := a.Insert(r.Context(), body.Collection, body.ID, body.Content)
	if err != nil {
		return http.StatusInternalServerError, errors.New("id[" + body.ID + "] insert handler failed: " + err.Error())
	}

	var response struct {
		Timestamp time.Time `json:"ts"`
	}
	response.Timestamp = insertedTime
	marshalledRsp, err := json.Marshal(response)
	if err != nil {
		a.Logger.Warn("failed to marshal time: "+err.Error(), zap.String("id", body.ID))
	}

	if _, err := w.Write(marshalledRsp); err != nil {
		a.Logger.Warn("failed to write the response: "+err.Error(), zap.String("id", body.ID), zap.String("reponse", string(marshalledRsp)))
	}

	return http.StatusOK, nil
}

func delete(a *app.App, w http.ResponseWriter, r *http.Request) (int, error) {

	var id string
	var collection string
	if err := multierr.Combine(getDocumentID(r, &id), getCollection(r, &collection)); err != nil {
		return http.StatusBadRequest, err
	}

	a.Logger.Debug("handling delete request", zap.String("id", id), zap.String("collection", collection))

	err := a.Delete(r.Context(), collection, id)
	if err != nil {
		return http.StatusInternalServerError, errors.New("id[" + id + "] delete handler failed: " + err.Error())
	}

	return http.StatusOK, nil
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
