package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
)

//API is the current API version
const API = "2.0"
const apiPath = "/api/" + API

//Router returns a new API router
func (s *Server) Router() http.Handler {
	r := mux.NewRouter()

	api := r.PathPrefix(apiPath).Subrouter()

	api.NotFoundHandler = withJSONResponse(func(r *http.Request) (int, interface{}) {
		return http.StatusNotFound, nil
	})

	api.Methods("POST").Path("/auth").Handler(
		withLogging("Authenticate", s.output,
			withJSONResponse(
				s.authenticate)))

	api.Methods("POST").Path("/auth/ping").Handler(
		withLogging("Ping", s.output,
			withJSONResponse(
				withAuth(s.sessionStore,
					ping))))

	api.Methods("GET").Path("/devices/{id:[0-9]{4}}").Handler(
		withLogging("ReadDevice", s.output,
			withJSONResponse(
				withAuth(s.sessionStore,
					withTX(s.db, s.readDevice)))))

	api.Methods("POST").Path("/devices/{id:[0-9]{4}}/checkin").Handler(
		withLogging("CheckinDevice", s.output,
			withJSONResponse(
				withAuth(s.sessionStore,
					withTX(s.db, s.checkinDevice)))))

	return r
}
