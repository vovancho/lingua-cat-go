package main

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
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

type customKafkaMarshaler struct{}

func (m customKafkaMarshaler) Marshal(topic string, msg *message.Message) (*sarama.ProducerMessage, error) {
	// Логируем перед отправкой
	logger.Info("Marshaling message for Kafka", map[string]interface{}{
		"topic":   topic,
		"uuid":    msg.UUID,
		"payload": string(msg.Payload),
	})

	// Создаем новое сообщение с тем же содержимым
	producerMsg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(msg.UUID),
		Value: sarama.ByteEncoder(msg.Payload),
	}

	// Добавляем метаданные
	for k, v := range msg.Metadata {
		producerMsg.Headers = append(producerMsg.Headers, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	return producerMsg, nil
}

func (m customKafkaMarshaler) Unmarshal(msg *sarama.ConsumerMessage) (*message.Message, error) {
	// Не используется в Publisher, но должно быть реализовано
	return message.NewMessage(string(msg.Key), msg.Value), nil
}

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

	// Настройка Sarama для отладки
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Net.MaxOpenRequests = 1

	saramaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:               []string{"lcg-kafka:9092"},
			Marshaler:             customKafkaMarshaler{},
			OverwriteSaramaConfig: saramaConfig,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	// Middleware для проверки сообщения перед отправкой
	forwarderMiddleware := func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			logger.Info("Forwarding message to Kafka", map[string]interface{}{
				"uuid":    msg.UUID,
				"payload": string(msg.Payload),
				"topic":   "lcg_exercise_topic",
			})
			return h(msg)
		}
	}

	config := forwarder.Config{
		ForwarderTopic: "lcg_exercise_topic",
		Middlewares:    []message.HandlerMiddleware{forwarderMiddleware},
		CloseTimeout:   30 * time.Second,
	}

	fwd, err := forwarder.NewForwarder(subscriber, saramaPublisher, logger, config)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := fwd.Run(context.Background()); err != nil {
			panic(err)
		}
	}()

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
		msg.Metadata.Set("published_at", time.Now().UTC().Format(time.RFC3339))

		logger.Info("Creating outbox message", map[string]interface{}{
			"uuid":    msg.UUID,
			"payload": string(msg.Payload),
		})

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

		if err := txPublisher.Publish("lcg_exercise_topic", msg); err != nil {
			return err
		}

		return tx.Commit()
	}

	time.Sleep(1 * time.Second)

	if err := insertAndPublish("b7237299-1e3a-4d46-bbdd-68f015598657", "en", 10); err != nil {
		panic(err)
	}

	select {}
}
