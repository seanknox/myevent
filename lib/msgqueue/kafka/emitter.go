package kafka

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/seanknox/myevent/lib/msgqueue"
)

type kafkaEventEmitter struct {
	producer sarama.SyncProducer
}

type messageEnvelope struct {
	EventName string      `json:"eventName"`
	Payload   interface{} `json:"payload"`
}

func NewKafkaEventEmitterFromEnvironment() (msgqueue.EventEmitter, error) {
	brokers := []string{"localhost:9092"}

	if brokerList := os.Getenv("KAFKA_BROKERS"); brokerList != "" {
		brokers = strings.Split(brokerList, ",")
	}

	config := sarama.NewConfig()
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}
	return NewKafkaEventEmitter(client)
}

func NewKafkaEventEmitter(client sarama.Client) (msgqueue.EventEmitter, error) {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	emitter := kafkaEventEmitter{
		producer: producer,
	}

	return &emitter, nil
}

func (k *kafkaEventEmitter) Emit(e msgqueue.Event) error {
	envelope := messageEnvelope{e.EventName(), e}
	jsonBody, err := json.Marshal(&envelope)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: e.EventName(),
		Value: sarama.ByteEncoder(jsonBody),
	}

	_, _, err = k.producer.SendMessage(msg)

	return err
}
