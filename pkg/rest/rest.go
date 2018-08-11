package rest

import (
	"net/http"

	"github.com/seanknox/myevent/pkg/persistence"

	"github.com/gorilla/mux"
)

func ServeAPI(endpoint string, dbHandler persistence.DatabaseHandler) error {
	handler := newEventHandler(dbHandler)
	r := mux.NewRouter()

	eventsrouter := r.PathPrefix("/events").Subrouter()

	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.allEventsHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)

	// return http.ListenAndServe(endpoint, r)
	return http.ListenAndServeTLS(endpoint, "cert.pem", "key.pem", r)

}
