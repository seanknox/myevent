package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/seanknox/myevent/lib/persistence/dblayer"
)

var (
	DBTypeDefault              = dblayer.DBTYPE("mongodb")
	DBConnectionDefault        = "mongodb://127.0.0.1:27017"
	RestfulEPDefault           = "localhost:8182"
	MessageBrokerTypeDefault   = "kafka"
	AMQPMessageBrokerDefault   = "amqp://guest:guest@localhost:5672"
	KafkaMessageBrokersDefault = []string{"localhost:9092"}
	ZipkinURIDefault           = "http://localhost:9411/api/v2/spans"
)

type ServiceConfig struct {
	DatabaseType        dblayer.DBTYPE `json:"databasetype"`
	DBConnection        string         `json:"dbconnection"`
	RestfulEndpoint     string         `json:"restfulapi_endpoint"`
	MessageBrokerType   string         `json:"message_broker_type"`
	AMQPMessageBroker   string         `json:"amqp_message_broker"`
	KafkaMessageBrokers []string       `json:"kafka_message_brokers"`
	ZipkinURI           string         `json:"zipkin_uri"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
		MessageBrokerTypeDefault,
		AMQPMessageBrokerDefault,
		KafkaMessageBrokersDefault,
		ZipkinURIDefault,
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values.")
		return conf, err
	}

	err = json.NewDecoder(file).Decode(&conf)

	if v := os.Getenv("AMQP_BROKER_URL"); v != "" {
		conf.MessageBrokerType = "amqp"
		conf.AMQPMessageBroker = v
	} else if v := os.Getenv("KAFKA_BROKER_URLS"); v != "" {
		conf.MessageBrokerType = "kafka"
		conf.KafkaMessageBrokers = strings.Split(v, ",")
	}
	return conf, err
}
