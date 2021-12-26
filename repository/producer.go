package repository

import (
	"fmt"
	"freq/config"
	"freq/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/vmihailenco/msgpack/v5"
)

func ProducerMessage(email *models.Email) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.Config("BOOTSTRAP_SERVER"),
		"security.protocol": "SASL_SSL",
		"sasl.username":     config.Config("USERNAME"),
		"sasl.password":     config.Config("PASSWORD"),
		"sasl.mechanism":    "PLAIN",
	})

	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	b, err := msgpack.Marshal(*email)

	if err != nil {
		panic(err)
	}
	// Produce messages to topic (asynchronously)
	topic := config.Config("PRODUCE_TOPIC")

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}