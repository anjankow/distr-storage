package server

import (
	"fmt"
	"net/http"
	"node/internal/app"
	"node/internal/config"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

type server struct {
	logger     *zap.Logger
	app        *app.App
	httpServer *http.Server
	addr       string
}

type appHandler struct {
	app    *app.App
	Handle AppHandleFunc
}

type AppHandleFunc func(*app.App, http.ResponseWriter, *http.Request) (int, error)

func (ser server) registerHandlers(router *mux.Router) {

	router.HandleFunc("/health", healthcheck)

	router.
		Path("/doc").
		Methods(http.MethodPut).
		Handler(appHandler{app: ser.app, Handle: insert})

	router.
		Path("/doc").
		Methods(http.MethodGet).
		Handler(appHandler{app: ser.app, Handle: get})
}

func NewServer(logger *zap.Logger, a *app.App) server {

	addr := config.GetPort()
	logger.Info(fmt.Sprint("listening on address: ", addr))

	return server{
		logger: logger,
		app:    a,
		addr:   addr,
	}
}

func (ser server) Run() error {
	router := mux.NewRouter()
	ser.registerHandlers(router)

	ser.httpServer = &http.Server{
		Handler:  router,
		ErrorLog: zap.NewStdLog(ser.logger),
		Addr:     ser.addr,
	}

	return ser.httpServer.ListenAndServe()
}

func (appHndl appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	status, err := appHndl.Handle(appHndl.app, w, r)

	if err != nil {
		appHndl.app.Logger.Warn("request failed", zap.Error(err))

		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}

		w.Write([]byte(fmt.Sprintln(err)))
		return
	}
}
