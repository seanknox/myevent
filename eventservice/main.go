package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/seanknox/myevent/lib/msgqueue"

	"github.com/seanknox/myevent/eventservice/pkg/rest"
	msgqueue_amqp "github.com/seanknox/myevent/lib/msgqueue/amqp"
	"github.com/seanknox/myevent/lib/msgqueue/kafka"
	"github.com/seanknox/myevent/lib/persistence/dblayer"
	"github.com/streadway/amqp"

	"github.com/seanknox/myevent/lib/config"
)

func main() {
	var emitter msgqueue.EventEmitter
	confPath := flag.String("config", `./config/config.json`, "flag to set path of configuration file")
	flag.Parse()

	// extract config
	config, _ := config.ExtractConfiguration(*confPath)

	switch config.MessageBrokerType {
	case "amqp":
		conn, err := amqp.Dial(config.AMQPMessageBroker)
		if err != nil {
			panic(err)
		}
		emitter, err = msgqueue_amqp.NewAMQPEventEmitter(conn)
		if err != nil {
			panic(err)
		}
	case "kafka":
		conf := sarama.NewConfig()
		conf.Producer.Return.Successes = true
		conn, err := sarama.NewClient(config.KafkaMessageBrokers, conf)
		if err != nil {
			panic(err)
		}

		emitter, err = kafka.NewKafkaEventEmitter(conn)
		if err != nil {
			panic(err)
		}
	default:
		panic("Bad message broker type: " + config.MessageBrokerType)
	}

	fmt.Println("Event Service connecting to database...")
	dbhandler, err := dblayer.NewPersistenceLayer(config.DatabaseType, config.DBConnection)
	if err != nil {
		log.Fatalf("Event Service couldn't connect to database: %+v", err)
	}
	fmt.Println("Event Service connected to DB.")

	// API start
	rest.ServeAPI(config.EventServiceEndpoint, config.ZipkinURI, dbhandler, emitter)
}
