package httpapi

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/korylprince/bisd-device-checkin-server/api"
)

//NewRouter returns an HTTP router for the HTTP API
func NewRouter(w io.Writer, config *api.AuthConfig, s SessionStore, db *sql.DB) http.Handler {

	//construct middleware
	var m = func(h returnHandler) http.Handler {
		return logMiddleware(jsonMiddleware(txMiddleware(authMiddleware(h, s), db)), w)
	}

	r := mux.NewRouter()

	r.Path("/devices/{bagTag:[0-9]{4}}").Methods("GET").Handler(m(handleReadDevice))
	r.Path("/devices/{bagTag:[0-9]{4}}/checkin").Methods("POST").Handler(m(handleCheckinDevice))

	r.Path("/auth").Methods("POST").Handler(logMiddleware(jsonMiddleware(txMiddleware(handleAuthenticate(config, s), db)), w))

	r.NotFoundHandler = m(notFoundHandler)

	return http.StripPrefix("/api/1.0", r)
}
