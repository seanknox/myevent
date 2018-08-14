package rest

import (
	"net/http"
	"time"

	"github.com/seanknox/myevent/lib/msgqueue"
	"github.com/seanknox/myevent/lib/persistence"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func ServeAPI(listenAddr string, database persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) {
	r := mux.NewRouter()
	r.Methods("POST").Path("/events/{eventID}/bookings").Handler(&CreateBookingHandler{eventEmitter, database})

	srv := http.Server{
		Handler:      handlers.CORS()(r),
		Addr:         listenAddr,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	srv.ListenAndServe()
}
