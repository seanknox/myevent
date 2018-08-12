package main

import (
	"flag"

	"github.com/seanknox/myevent/bookingservice/pkg/listener"
	"github.com/seanknox/myevent/bookingservice/pkg/rest"
	"github.com/seanknox/myevent/lib/config"
	"github.com/seanknox/myevent/lib/msgqueue"
	msgqueue_amqp "github.com/seanknox/myevent/lib/msgqueue/amqp"
	"github.com/seanknox/myevent/lib/persistence/dblayer"
	"github.com/streadway/amqp"
)

func main() {
	var eventListener msgqueue.EventListener
	var eventEmitter msgqueue.EventEmitter

	confPath := flag.String("config", `./config/config.json`, "flag to set path of configuration file")
	flag.Parse()

	config, _ := config.ExtractConfiguration(*confPath)

	connection, err := amqp.Dial(config.AMQPMessageBroker)
	if err != nil {
		panic("could not establish AMQP connection: " + err.Error())
	}

	eventListener, err = msgqueue_amqp.NewAMQPEventListener(connection, "events")
	if err != nil {
		panic(err)
	}

	eventEmitter, err = msgqueue_amqp.NewAMQPEventEmitter(connection)
	if err != nil {
		panic(err)
	}

	dbhandler, err := dblayer.NewPersistenceLayer(config.DatabaseType, config.DBConnection)
	if err != nil {
		panic("could not establish database connection: " + err.Error())
	}

	processor := listener.EventProcessor{eventListener, dbhandler}
	go processor.ProcessEvents()

	rest.ServeAPI(config.RestfulEndpoint, dbhandler, eventEmitter)

}
