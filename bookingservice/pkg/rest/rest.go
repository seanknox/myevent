package rest

import (
	"log"
	"net/http"
	"time"

	"github.com/seanknox/myevent/lib/msgqueue"
	"github.com/seanknox/myevent/lib/persistence"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

func ServeAPI(listenAddr, zipkin_uri string, database persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) {
	// setup zipkin span reporter
	reporter := httpreporter.NewReporter(zipkin_uri)
	defer reporter.Close()

	// create local zipkin endpoint
	zipkinEndpoint, err := zipkin.NewEndpoint("booking", listenAddr)
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
	r := mux.NewRouter()
	r.Use(serverMiddleware)
	r.Methods("POST").Path("/events/{eventID}/book").Handler(&CreateBookingHandler{eventEmitter, database})

	srv := http.Server{
		Handler:      handlers.CORS()(r),
		Addr:         listenAddr,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	srv.ListenAndServe()
}
