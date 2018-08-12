package rest

import (
	"net/http"

	"github.com/seanknox/myevent/lib/msgqueue"

	"github.com/seanknox/myevent/lib/persistence"

	"github.com/gorilla/mux"
)

func ServeAPI(endpoint, tlsendpoint string, dbHandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) (chan error, chan error) {
	handler := newEventHandler(dbHandler, eventEmitter)
	r := mux.NewRouter()

	eventsrouter := r.PathPrefix("/events").Subrouter()

	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.allEventsHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)

	httpErrChan := make(chan error)
	httpsErrChan := make(chan error)

	go func() {
		httpErrChan <- http.ListenAndServe(endpoint, r)
	}()

	go func() {
		httpsErrChan <- http.ListenAndServeTLS(tlsendpoint, "cert.pem", "key.pem", r)
	}()

	return httpErrChan, httpsErrChan
}
