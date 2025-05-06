package main

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

var logger = watermill.NewStdLogger(false, false)

func main() {
	db, err := sqlx.Connect("postgres", "postgres://exercise:secret@lcg-exercise-db:5432/exercise?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	subscriber, err := sql.NewSubscriber(
		db,
		sql.SubscriberConfig{
			SchemaAdapter:    sql.DefaultPostgreSQLSchema{},
			OffsetsAdapter:   sql.DefaultPostgreSQLOffsetsAdapter{},
			InitializeSchema: true,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	saramaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{"lcg-kafka:9092"},
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	//router, err := message.NewRouter(message.RouterConfig{}, logger)
	//if err != nil {
	//	panic(err)
	//}
	//
	//router.AddNoPublisherHandler(
	//	"outbox_to_kafka",
	//	"lcg_exercise_topic",
	//	subscriber,
	//	func(msg *message.Message) error {
	//		return saramaPublisher.Publish("lcg_exercise_topic", msg)
	//	},
	//)

	// Configure the forwarder
	config := forwarder.Config{
		ForwarderTopic:      "lcg_exercise_topic",
		Middlewares:         nil,
		CloseTimeout:        30 * time.Second,
		AckWhenCannotUnwrap: false,
		Router:              nil,
	}

	// Create the forwarder
	fwd, err := forwarder.NewForwarder(subscriber, saramaPublisher, logger, config)
	if err != nil {
		panic(err)
	}

	// Run the forwarder in a separate goroutine
	go func() {
		if err := fwd.Run(context.Background()); err != nil {
			panic(err)
		}
	}()

	// Function to insert into exercise table and publish to outbox in one transaction
	insertAndPublish := func(userID, lang string, taskAmount int) error {
		tx, err := db.Beginx()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		var id int
		err = tx.QueryRowx("INSERT INTO exercise (user_id, lang, task_amount) VALUES ($1, $2, $3) RETURNING id", userID, lang, taskAmount).Scan(&id)
		if err != nil {
			return err
		}

		payload, err := json.Marshal(map[string]interface{}{
			"id":                id,
			"user_id":           userID,
			"lang":              lang,
			"task_amount":       taskAmount,
			"destination_topic": "lcg_exercise_topic",
		})
		if err != nil {
			return err
		}

		msg := message.NewMessage(watermill.NewUUID(), payload)

		// Use the transaction for publishing
		txPublisher, err := sql.NewPublisher(
			tx,
			sql.PublisherConfig{
				SchemaAdapter: sql.DefaultPostgreSQLSchema{},
			},
			logger,
		)
		if err != nil {
			return err
		}

		err = txPublisher.Publish("lcg_exercise_topic", msg)
		if err != nil {
			return err
		}

		// Simulate error for testing
		//time.Sleep(1 * time.Second)
		//panic("test error")

		return tx.Commit()
	}

	time.Sleep(1 * time.Second)

	// Example usage
	if err := insertAndPublish("b7237299-1e3a-4d46-bbdd-68f015598657", "en", 10); err != nil {
		panic(err)
	}

	// Infinite loop to keep the application running
	select {}
}
