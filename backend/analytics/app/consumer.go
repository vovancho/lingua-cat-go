package main

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"log"
)

func main() {
	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:       []string{"lcg-kafka:9092"},
			ConsumerGroup: "exercise_completed_group101",
			Unmarshaler:   kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()

	messages, err := subscriber.Subscribe(context.Background(), "lcg_exercise_completed")
	if err != nil {
		panic(err)
	}

	go func() {
		for msg := range messages {
			log.Printf("Received message: %s", msg.UUID)
			log.Println(msg.Payload)
			msg.Ack()
		}
	}()

	select {}
}
