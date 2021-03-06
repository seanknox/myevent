package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/seanknox/myevent/contracts"
	"github.com/seanknox/myevent/lib/msgqueue"
	"github.com/seanknox/myevent/lib/persistence"

	"github.com/gorilla/mux"
)

type eventServiceHandler struct {
	dbHandler    persistence.DatabaseHandler
	eventEmitter msgqueue.EventEmitter
}

func newEventHandler(databaseHandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) *eventServiceHandler {
	return &eventServiceHandler{
		dbHandler:    databaseHandler,
		eventEmitter: eventEmitter,
	}
}

func (eh *eventServiceHandler) findEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	criteria, ok := vars["SearchCriteria"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error": "no search criteria found. You can search by either id via /id/[num] or by name via /name/[name]"}`)
		return
	}

	searchKey, ok := vars["search"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error": "no search keys found, e.g. /id/4 or /name/badreligionconcert"}`)
		return
	}

	var event persistence.Event
	var err error
	switch strings.ToLower(criteria) {
	case "name":
		event, err = eh.dbHandler.FindEventByName(searchKey)
	case "id":
		id, err := hex.DecodeString(searchKey)
		if err == nil {
			event, err = eh.dbHandler.FindEvent(id)
		}
	}

	if err != nil {
		fmt.Fprintf(w, `{"error": "%s"}`, err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&event)
}

func (eh *eventServiceHandler) allEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := eh.dbHandler.FindAllAvailableEvents()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"error": "Error occured while trying to find all available events %s"}`, err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err = json.NewEncoder(w).Encode(&events)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"error": "Error occured while trying encode events to JSON %s"}`, err)
	}

}

func (eh *eventServiceHandler) newEventHandler(w http.ResponseWriter, r *http.Request) {
	event := persistence.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"error": "error occured while decoding event data %s"}`, err)
		return
	}
	id, err := eh.dbHandler.AddEvent(event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"error": "error occured while persisting event %d %s"}`, id, err)
		return
	}

	msg := contracts.EventCreatedEvent{
		ID:         hex.EncodeToString(id),
		Name:       event.Name,
		LocationID: string(event.Location.ID),
		Start:      time.Unix(event.StartDate, 0),
		End:        time.Unix(event.EndDate, 0),
	}

	eh.eventEmitter.Emit(&msg)
	fmt.Fprint(w, `{"id": "%d"}`, id)
}
