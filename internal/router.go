package internal

import (
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/gorilla/mux"
	"net/http"
)

func loadRoutes(r *mux.Router, _ *di.Container) {
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/ping", ping).Methods(http.MethodGet)
}

func ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("pong"))
}
