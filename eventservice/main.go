package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/seanknox/myevent/eventservice/pkg/rest"
	msgqueue_amqp "github.com/seanknox/myevent/lib/msgqueue/amqp"
	"github.com/seanknox/myevent/lib/persistence/dblayer"
	"github.com/streadway/amqp"

	"github.com/seanknox/myevent/lib/config"
)

func main() {
	confPath := flag.String("config", `./config/config.json`, "flag to set path of configuration file")
	flag.Parse()

	// extract config
	config, _ := config.ExtractConfiguration(*confPath)

	conn, err := amqp.Dial(config.AMQPMessageBroker)
	if err != nil {
		panic(err)
	}
	emitter, err := msgqueue_amqp.NewAMQPEventEmitter(conn)
	if err != nil {
		panic(err)
	}

	fmt.Println("Event Service connecting to database...")
	dbhandler, err := dblayer.NewPersistenceLayer(config.DatabaseType, config.DBConnection)
	if err != nil {
		log.Fatalf("Event Service couldn't connect to database: %+v", err)
	}
	fmt.Println("Event Service connected to DB.")

	// API start
	httpErrChan, httpsErrChan := rest.ServeAPI(config.RestfulEndpoint, config.RestfulTLSEndpoint, dbhandler, emitter)

	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP error: ", err)
	case err := <-httpsErrChan:
		log.Fatal("HTTPS error: ", err)
	}
}
