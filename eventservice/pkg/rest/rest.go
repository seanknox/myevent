package rest

import (
	"log"
	"net/http"

	"github.com/seanknox/myevent/lib/msgqueue"

	"github.com/seanknox/myevent/lib/persistence"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

func ServeAPI(endpoint, zipkin_uri string, dbHandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) error {
	// setup zipkin span reporter
	reporter := httpreporter.NewReporter(zipkin_uri)
	defer reporter.Close()

	// create local zipkin endpoint
	zipkinEndpoint, err := zipkin.NewEndpoint("event", endpoint)
	if err != nil {
		log.Fatalf("unable to create local zipking endpoint: %+v\n", err)
	}

	// initialize tracer
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zipkinEndpoint))
	if err != nil {
		log.Fatalf("unable to initialize tracer %+v\n", err)
	}

	// create global zipkin http server middleware
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		tracer, zipkinhttp.TagResponseSize(true),
	)

	handler := newEventHandler(dbHandler, eventEmitter)
	r := mux.NewRouter()
	r.Use(serverMiddleware)
	server := handlers.CORS()(r)

	eventsrouter := r.PathPrefix("/events").Subrouter()

	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.allEventsHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)

	return http.ListenAndServe(endpoint, server)

}
