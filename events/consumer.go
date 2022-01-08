package events

import (
	"fmt"
	"freq/config"
	"freq/models"
	"freq/repository"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/vmihailenco/msgpack/v5"
)

func CreateConsumer() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka1:9092",
		"group.id":          config.Config("CONSUMER_GROUP_ID"),
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	err = c.SubscribeTopics([]string{config.Config("CONSUME_TOPIC")}, nil)
	if err != nil {
		panic(err)
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			email := new(models.Email)
			err = msgpack.Unmarshal(msg.Value, email)
			fmt.Printf("Message on %s: %v\n", msg.TopicPartition, email)
			fmt.Println(email)
			_ = repository.EmailRepoImpl{}.UpdateEmailStatus(email.Id, email.Status)
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}
