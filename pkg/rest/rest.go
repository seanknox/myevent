package rest

import (
	"log"
	"net/http"

	"github.com/seanknox/myevent/pkg/persistence"

	"github.com/gorilla/mux"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

func ServeAPI(endpoint string, dbHandler persistence.DatabaseHandler) error {
	// setup zipkin span reporter
	reporter := httpreporter.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	// create local zipkin endpoint
	zipkinEndpoint, err := zipkin.NewEndpoint("myEvent", endpoint)
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

	handler := newEventHandler(dbHandler)
	r := mux.NewRouter()
	r.Use(serverMiddleware)

	eventsrouter := r.PathPrefix("/events").Subrouter()

	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.allEventsHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)

	// return http.ListenAndServe(endpoint, r)
	return http.ListenAndServeTLS(endpoint, "cert.pem", "key.pem", r)

}
